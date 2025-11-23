package sqlite

import (
	"context"
	"database/sql"

	"pr-review/internal/errors"
	"pr-review/internal/models"
)

func (r *SQLiteRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	const op = "SQLite.CreateTeam"

	exists, err := r.TeamExists(ctx, team.Name)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if exists {
		return errors.WrapError(op, errors.ErrTeamExists)
	}

	for _, member := range team.Members {
		exists, err := r.UserExists(ctx, member.UserID, member.Username)
		if err != nil {
			return errors.WrapError(op, err)
		}
		if exists {
			return errors.WrapError(op, errors.ErrUserExists)
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.WrapError(op, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			return
		}
	}()

	query := `INSERT INTO teams (name) VALUES (?)`
	_, err = tx.ExecContext(ctx, query, team.Name)
	if err != nil {
		return errors.WrapError(op, err)
	}

	userQuery := `INSERT INTO users (user_id, username, is_active, team_name) VALUES (?, ?, ?, ?)`
	for _, member := range team.Members {
		_, err := tx.ExecContext(ctx, userQuery, member.UserID, member.Username, member.IsActive, team.Name)
		if err != nil {
			return errors.WrapError(op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.WrapError(op, err)
	}
	return nil
}

func (r *SQLiteRepository) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
	const op = "SQLite.GetTeamByName"

	exists, err := r.TeamExists(ctx, name)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	if !exists {
		return nil, errors.WrapError(op, errors.ErrTeamNotFound)
	}

	query := `
		SELECT user_id, username, is_active 
		FROM users 
		WHERE team_name = ? 
		ORDER BY username
	`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, errors.WrapError(op, err)
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WrapError(op, err)
	}

	team := &models.Team{
		Name:    name,
		Members: members,
	}

	return team, nil
}

func (r *SQLiteRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	const op = "SQLite.TeamExists"

	query := `SELECT 1 FROM teams WHERE name = ?`
	row := r.db.QueryRowContext(ctx, query, teamName)

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

func (r *SQLiteRepository) GetPRsCntByTeam(ctx context.Context, teamName string) (int, error) {
	const op = "SQLite.GetPRsCntByTeam"

	exists, err := r.TeamExists(ctx, teamName)
	if err != nil {
		return 0, errors.WrapError(op, err)
	}
	if !exists {
		return 0, errors.WrapError(op, errors.ErrTeamNotFound)
	}

	query := `
		SELECT COUNT(pr.id) as prs_authored
		FROM users u
		LEFT JOIN pull_requests pr ON u.user_id = pr.author_id
		WHERE u.team_name = ?
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, teamName)

	var count int
	err = row.Scan(&count)
	if err == sql.ErrNoRows {
		return 0, errors.WrapError(op, errors.ErrTeamNotFound)
	}
	if err != nil {
		return 0, errors.WrapError(op, err)
	}

	return count, nil
}

func (r *SQLiteRepository) GetAvgReviewersPerPR(ctx context.Context, teamName string) (float64, error) {
	const op = "SQLite.GetPRsCntByTeam"

	exists, err := r.TeamExists(ctx, teamName)
	if err != nil {
		return 0, errors.WrapError(op, err)
	}
	if !exists {
		return 0, errors.WrapError(op, errors.ErrTeamNotFound)
	}

	query := `
		SELECT 
			ROUND(
				CAST(COUNT(prr.user_id) AS FLOAT) / 
				NULLIF(COUNT(DISTINCT pr.id), 0), 
				2
			) as avg_reviewers_per_pr
		FROM users u
		LEFT JOIN pull_requests pr ON u.user_id = pr.author_id
		LEFT JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE u.team_name = ?
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, teamName)

	var count float64
	err = row.Scan(&count)
	if err == sql.ErrNoRows {
		return 0, errors.WrapError(op, errors.ErrTeamNotFound)
	}
	if err != nil {
		return 0, errors.WrapError(op, err)
	}

	return count, nil
}
