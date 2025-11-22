package sqlite

import (
	"context"
	"database/sql"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/service"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) service.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const op = "SQLite.GetUserByID"

	exists, err := r.UserExists(ctx, id)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `SELECT user_id, username, is_active, team_name FROM users WHERE user_id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err = row.Scan(&user.UserID, &user.Username, &user.IsActive, &user.TeamName)
	if err == sql.ErrNoRows {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	const op = "SQLite.GetUserByUsername"

	query := `SELECT user_id, username, is_active, team_name FROM users WHERE username = ?`
	row := r.db.QueryRowContext(ctx, query, username)

	var user models.User
	err := row.Scan(&user.UserID, &user.Username, &user.IsActive, &user.TeamName)
	if err == sql.ErrNoRows {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	return &user, nil
}

func (r *userRepository) SetUserActive(ctx context.Context, userID string, isActive bool) error {
	const op = "SQLite.SetUserActive"

	exists, err := r.UserExists(ctx, userID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if !exists {
		return errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `UPDATE users SET is_active = ? WHERE user_id = ?`
	result, err := r.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return errors.WrapError(op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(op, err)
	}
	if rows == 0 {
		return errors.WrapError(op, errors.ErrUserNotFound)
	}

	return nil
}

func (r *userRepository) GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	const op = "SQLite.GetPRsByReviewer"

	exists, err := r.UserExists(ctx, userID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE prr.user_id = ? AND pr.status = 'OPEN'
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer rows.Close()

	var prs []*models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status)
		if err != nil {
			return nil, errors.WrapError(op, err)
		}
		prs = append(prs, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WrapError(op, err)
	}

	return prs, nil
}

func (r *userRepository) UserExists(ctx context.Context, userID string) (bool, error) {
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
