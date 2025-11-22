package sqlite

import (
	"context"
	"database/sql"
	"math/rand"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/service"
	"time"
)

type prRepository struct {
	db *sql.DB
}

func NewprRepository(db *sql.DB) service.PRRepository {
	return &prRepository{db: db}
}

func (r *prRepository) CreatePR(ctx context.Context, pr *models.PullRequestShort) error {
	const op = "SQLite.CreatePR"

	exists, err := r.PRExists(ctx, pr.ID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if exists {
		return errors.WrapError(op, errors.ErrPRExists)
	}

	exists, err = r.userExists(ctx, pr.AuthorID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if !exists {
		return errors.WrapError(op, errors.ErrUserNotFound)
	}

	authorTeam, err := r.getUserTeam(ctx, pr.AuthorID)
	if err != nil {
		return errors.WrapError(op, err)
	}

	availableReviewers, err := r.getActiveTeamMembers(ctx, authorTeam, pr.AuthorID, nil)
	if err != nil {
		return errors.WrapError(op, err)
	}

	selectedReviewers := r.selectRandomReviewers(availableReviewers, 2)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.WrapError(op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO pull_requests (id, name, author_id, status, created_at) VALUES (?, ?, ?, 'OPEN', ?)`
	_, err = tx.ExecContext(ctx, query, pr.ID, pr.Name, pr.AuthorID, time.Now())
	if err != nil {
		return errors.WrapError(op, err)
	}

	if len(selectedReviewers) > 0 {
		reviewersQuery := `INSERT INTO pr_reviewers (pr_id, user_id) VALUES (?, ?)`
		for _, reviewerID := range selectedReviewers {
			_, err := tx.ExecContext(ctx, reviewersQuery, pr.ID, reviewerID)
			if err != nil {
				return errors.WrapError(op, err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WrapError(op, err)
	}

	return nil
}

func (r *prRepository) GetPRByID(ctx context.Context, id string) (*models.PullRequest, error) {
	const op = "SQLite.GetPRByID"

	exists, err := r.PRExists(ctx, id)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrPRNotFound)
	}

	query := `SELECT id, name, author_id, status, created_at, merged_at FROM pull_requests WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var pr models.PullRequest
	var mergedAt sql.NullTime

	err = row.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)
	if err == sql.ErrNoRows {
		return nil, errors.WrapError(op, errors.ErrPRNotFound)
	}
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	reviewers, err := r.getPRReviewers(ctx, id)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *prRepository) MergePR(ctx context.Context, prID string, mergedAt time.Time) error {
	const op = "SQLite.MergePR"

	exists, err := r.PRExists(ctx, prID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if !exists {
		return errors.WrapError(op, errors.ErrPRNotFound)
	}

	query := `UPDATE pull_requests SET status = 'MERGED', merged_at = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, mergedAt, prID)
	if err != nil {
		return errors.WrapError(op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(op, err)
	}
	if rows == 0 {
		return errors.WrapError(op, errors.ErrPRNotFound)
	}

	return nil
}

func (r *prRepository) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*string, error) {
	const op = "SQLite.ReassignReviewer"

	pr, err := r.GetPRByID(ctx, prID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if pr.Status == "MERGED" {
		return nil, errors.WrapError(op, errors.ErrPRMerged)
	}

	exists, err := r.userExists(ctx, oldUserID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}

	isAssigned, err := r.IsReviewerAssigned(ctx, prID, oldUserID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !isAssigned {
		return nil, errors.WrapError(op, errors.ErrNotAssigned)
	}

	authorTeam, err := r.getUserTeam(ctx, pr.AuthorID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	currentReviewers := pr.AssignedReviewers

	availableReviewers, err := r.getActiveTeamMembers(ctx, authorTeam, oldUserID, currentReviewers)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	if len(availableReviewers) == 0 {
		return nil, errors.WrapError(op, errors.ErrNoCandidate)
	}

	selectedReviewers := r.selectRandomReviewers(availableReviewers, 1)
	if len(selectedReviewers) == 0 {
		return nil, errors.WrapError(op, errors.ErrNoCandidate)
	}
	reviewerID := selectedReviewers[0]

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM pr_reviewers WHERE pr_id = ? AND user_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, prID, oldUserID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	insertQuery := `INSERT INTO pr_reviewers (pr_id, user_id) VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, insertQuery, prID, reviewerID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.WrapError(op, err)
	}

	return &reviewerID, nil
}

func (r *prRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	const op = "SQLite.PRExists"

	query := `SELECT 1 FROM pull_requests WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, prID)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.WrapError(op, err)
	}

	return true, nil
}

func (r *prRepository) IsReviewerAssigned(ctx context.Context, prID, userID string) (bool, error) {
	const op = "SQLite.IsReviewerAssigned"

	query := `SELECT 1 FROM pr_reviewers WHERE pr_id = ? AND user_id = ?`
	row := r.db.QueryRowContext(ctx, query, prID, userID)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.WrapError(op, err)
	}

	return true, nil
}

// private methods

func (r *prRepository) getPRReviewers(ctx context.Context, prID string) ([]string, error) {
	const op = "SQLite.getPRReviewers"

	query := `SELECT user_id FROM pr_reviewers WHERE pr_id = ? ORDER BY user_id`
	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, errors.WrapError(op, err)
		}
		reviewers = append(reviewers, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WrapError(op, err)
	}

	return reviewers, nil
}

func (r *prRepository) getUserTeam(ctx context.Context, userID string) (string, error) {
	const op = "SQLite.getUserTeam"

	query := `SELECT team_name FROM users WHERE user_id = ?`
	row := r.db.QueryRowContext(ctx, query, userID)

	var teamName string
	err := row.Scan(&teamName)
	if err == sql.ErrNoRows {
		return "", errors.ErrUserNotFound
	}
	if err != nil {
		return "", errors.WrapError(op, err)
	}

	return teamName, nil
}

func (r *prRepository) userExists(ctx context.Context, userID string) (bool, error) {
	const op = "SQLite.UserExists"

	query := `SELECT 1 FROM users WHERE user_id = ?`
	row := r.db.QueryRowContext(ctx, query, userID)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errors.WrapError(op, err)
	}

	return true, nil
}

func (r *prRepository) getActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string, excludeReviewers []string) ([]string, error) {
	const op = "SQLite.getActiveTeamMembers"

	excludeMap := make(map[string]bool)
	for _, reviewer := range excludeReviewers {
		excludeMap[reviewer] = true
	}

	query := `
		SELECT user_id 
		FROM users 
		WHERE team_name = ? AND is_active = TRUE AND user_id != ?
		ORDER BY user_id
	`
	rows, err := r.db.QueryContext(ctx, query, teamName, excludeUserID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			return nil, errors.WrapError(op, err)
		}
		if !excludeMap[userID] {
			members = append(members, userID)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WrapError(op, err)
	}

	return members, nil
}

func (r *prRepository) selectRandomReviewers(reviewers []string, maxCount int) []string {
	if len(reviewers) == 0 {
		return nil
	}

	if len(reviewers) <= maxCount {
		return reviewers
	}

	shuffled := make([]string, len(reviewers))
	copy(shuffled, reviewers)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:maxCount]
}
