package postgres

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	"pr-review/internal/errors"
	"pr-review/internal/models"
)

func (r *PostgresRepository) CreatePR(ctx context.Context, pr *models.PullRequestShort) error {
	const op = "Postgres.CreatePR"

	exists, err := r.PRExists(ctx, pr.ID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if exists {
		return errors.WrapError(op, errors.ErrPRExists)
	}

	exists, err = r.UserExists(ctx, pr.AuthorID)
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

	selectedReviewers := selectRandomReviewers(availableReviewers, 2)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.WrapError(op, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			return
		}
	}()

	query := `INSERT INTO pull_requests (id, name, author_id, status, created_at) VALUES ($1, $2, $3, 'OPEN', $4)`
	_, err = tx.ExecContext(ctx, query, pr.ID, pr.Name, pr.AuthorID, time.Now())
	if err != nil {
		return errors.WrapError(op, err)
	}

	if len(selectedReviewers) > 0 {
		reviewersQuery := `INSERT INTO pr_reviewers (pr_id, user_id) VALUES ($1, $2)`
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

func (r *PostgresRepository) GetPRByID(ctx context.Context, id string) (*models.PullRequest, error) {
	const op = "Postgres.GetPRByID"

	exists, err := r.PRExists(ctx, id)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrPRNotFound)
	}

	query := `SELECT id, name, author_id, status, created_at, merged_at FROM pull_requests WHERE id = $1`
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

func (r *PostgresRepository) MergePR(ctx context.Context, prID string, mergedAt time.Time) error {
	const op = "Postgres.MergePR"

	exists, err := r.PRExists(ctx, prID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if !exists {
		return errors.WrapError(op, errors.ErrPRNotFound)
	}

	query := `UPDATE pull_requests SET status = 'MERGED', merged_at = $1 WHERE id = $2 AND status = 'OPEN'`
	result, err := r.db.ExecContext(ctx, query, mergedAt, prID)
	if err != nil {
		return errors.WrapError(op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(op, err)
	}
	if rowsAffected == 0 {
		return errors.WrapError(op, errors.ErrPRNotFound)
	}

	return nil
}

func (r *PostgresRepository) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*string, error) {
	const op = "Postgres.ReassignReviewer"

	pr, err := r.GetPRByID(ctx, prID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if pr.Status == "MERGED" {
		return nil, errors.WrapError(op, errors.ErrPRMerged)
	}

	exists, err := r.UserExists(ctx, oldUserID)
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

	selectedReviewers := selectRandomReviewers(availableReviewers, 1)
	if len(selectedReviewers) == 0 {
		return nil, errors.WrapError(op, errors.ErrNoCandidate)
	}
	reviewerID := selectedReviewers[0]

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			return
		}
	}()

	deleteQuery := `DELETE FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2`
	_, err = tx.ExecContext(ctx, deleteQuery, prID, oldUserID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	insertQuery := `INSERT INTO pr_reviewers (pr_id, user_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, insertQuery, prID, reviewerID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.WrapError(op, err)
	}

	return &reviewerID, nil
}

func (r *PostgresRepository) PRExists(ctx context.Context, prID string) (bool, error) {
	const op = "Postgres.PRExists"

	query := `SELECT 1 FROM pull_requests WHERE id = $1`
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

func (r *PostgresRepository) IsReviewerAssigned(ctx context.Context, prID, userID string) (bool, error) {
	const op = "Postgres.IsReviewerAssigned"

	query := `SELECT 1 FROM pr_reviewers WHERE pr_id = $1 AND user_id = $2`
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

func (r *PostgresRepository) GetTotalStats(ctx context.Context) (*models.TotalStats, error) {
	const op = "Postgres.GetTotalStats"

	query := `
		SELECT 
			(SELECT COUNT(*) FROM teams) as total_teams,
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM users WHERE is_active = true) as active_users,
			(SELECT COUNT(*) FROM pull_requests) as total_prs,
			(SELECT COUNT(*) FROM pull_requests WHERE status = 'OPEN') as open_prs,
			(SELECT COUNT(*) FROM pull_requests WHERE status = 'MERGED') as merged_prs,
			COALESCE(
				ROUND(
					CAST((SELECT COUNT(*) FROM pr_reviewers) AS NUMERIC) / 
					NULLIF((SELECT COUNT(*) FROM pull_requests), 0), 
					2
				), 
				0
			) as avg_reviewers_per_pr
	`
	row := r.db.QueryRowContext(ctx, query)

	stats := &models.TotalStats{}
	err := row.Scan(
		&stats.TotalTeams,
		&stats.TotalUsers,
		&stats.ActiveUsers,
		&stats.TotalPRs,
		&stats.OpenPRs,
		&stats.MergedPRs,
		&stats.AvgReviewersPerPR,
	)
	if err == sql.ErrNoRows {
		return &models.TotalStats{}, nil
	}
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	return stats, nil
}

// private methods
func (r *PostgresRepository) getPRReviewers(ctx context.Context, prID string) ([]string, error) {
	const op = "Postgres.getPRReviewers"

	query := `SELECT user_id FROM pr_reviewers WHERE pr_id = $1 ORDER BY user_id`
	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

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

func (r *PostgresRepository) getUserTeam(ctx context.Context, userID string) (string, error) {
	const op = "Postgres.getUserTeam"

	query := `SELECT team_name FROM users WHERE user_id = $1`
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

func (r *PostgresRepository) getActiveTeamMembers(ctx context.Context, teamName, excludeUserID string, excludeReviewers []string) ([]string, error) {
	const op = "Postgres.getActiveTeamMembers"

	excludeMap := make(map[string]bool)
	for _, reviewer := range excludeReviewers {
		excludeMap[reviewer] = true
	}

	query := `
		SELECT user_id 
		FROM users 
		WHERE team_name = $1 AND is_active = TRUE AND user_id != $2
		ORDER BY user_id
	`
	rows, err := r.db.QueryContext(ctx, query, teamName, excludeUserID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

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

func selectRandomReviewers(reviewers []string, maxCount int) []string {
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
