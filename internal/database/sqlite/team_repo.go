package sqlite

import (
	"context"
	"database/sql"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/service"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(s *SQLiteDatabase) service.TeamRepository {
	return &teamRepository{db: s.db}
}

func (r *teamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	const op = "SQLite.CreateTeam"

	for _, member := range team.Members {
		exists, err := r.userExists(ctx, member.UserID, member.Username)
		if err != nil {
			return errors.WrapError(op, err)
		}
		if exists {
			return errors.WrapError(op, errors.ErrUserExists)
		}
	}

	exists, err := r.TeamExists(ctx, team.Name)
	if err != nil {
		return errors.WrapError(op, err)
	}
	if exists {
		return errors.WrapError(op, errors.ErrTeamExists)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.WrapError(op, err)
	}
	defer tx.Rollback()

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

func (r *teamRepository) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
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
	defer rows.Close()

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

func (r *teamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
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

// private methods

func (r *teamRepository) userExists(ctx context.Context, userID, username string) (bool, error) {
	const op = "SQLite.UserExists"

	query := `SELECT 1 FROM users WHERE user_id = ? OR username = ?`
	row := r.db.QueryRowContext(ctx, query, userID, username)

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
