package postgres

import (
	"context"
	"database/sql"

	"pr-review/internal/errors"
	"pr-review/internal/models"
)

func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const op = "Postgres.GetUserByID"

	exists, err := r.UserExists(ctx, id)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `SELECT user_id, username, is_active, team_name FROM users WHERE user_id = $1`
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

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	const op = "Postgres.GetUserByUsername"

	query := `SELECT user_id, username, is_active, team_name FROM users WHERE username = $1`
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

func (r *PostgresRepository) SetUserActive(ctx context.Context, userID string, isActive bool) error {
	const op = "Postgres.SetUserActive"

	exists, err := r.UserExists(ctx, userID)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if !exists {
		return errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `UPDATE users SET is_active = $1 WHERE user_id = $2`
	result, err := r.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return errors.WrapError(op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WrapError(op, err)
	}
	if rowsAffected == 0 {
		return errors.WrapError(op, errors.ErrUserNotFound)
	}

	return nil
}

func (r *PostgresRepository) GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	const op = "Postgres.GetPRsByReviewer"

	exists, err := r.UserExists(ctx, userID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status, pr.created_at
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE prr.user_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	var prs []*models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		var createdAt string
		err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &createdAt)
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

func (r *PostgresRepository) GetPRsCntByAuthor(ctx context.Context, userID string) (int, error) {
	const op = "Postgres.GetPRsCntByAuthor"

	exists, err := r.UserExists(ctx, userID)
	if err != nil {
		return 0, errors.WrapError(op, err)
	}
	if !exists {
		return 0, errors.WrapError(op, errors.ErrUserNotFound)
	}

	query := `
		SELECT COUNT(*)
		FROM pull_requests 
		WHERE author_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, userID)

	var count int
	err = row.Scan(&count)
	if err == sql.ErrNoRows {
		return 0, errors.WrapError(op, errors.ErrUserNotFound)
	}
	if err != nil {
		return 0, errors.WrapError(op, err)
	}

	return count, nil
}

func (r *PostgresRepository) UserExists(ctx context.Context, userID string, username ...string) (bool, error) {
	const op = "Postgres.UserExists"

	var row *sql.Row

	if len(username) > 0 {
		query := `SELECT 1 FROM users WHERE user_id = $1 OR username = $2`
		row = r.db.QueryRowContext(ctx, query, userID, username[0])
	} else {
		query := `SELECT 1 FROM users WHERE user_id = $1`
		row = r.db.QueryRowContext(ctx, query, userID)
	}

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
