package sqlite

import (
	"context"
	"database/sql"
	"pr-review/internal/database"
	"pr-review/internal/database/models"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) error {
	const op = "CreateTeam"

	exists, err := r.TeamExists(ctx, team.Name)
	if err != nil {
		return database.WrapError(op, err)
	}
	if exists {
		return database.WrapError(op, database.ErrTeamExists)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return database.WrapError(op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO teams (name) VALUES (?)`
	_, err = tx.ExecContext(ctx, query, team.Name)
	if err != nil {
		return database.WrapError(op, err)
	}

	userQuery := `INSERT INTO users (user_id, username, is_active, team_name) VALUES (?, ?, ?, ?)`
	for _, member := range team.Members {
		_, err := tx.ExecContext(ctx, userQuery, member.UserID, member.Username, member.IsActive, team.Name)
		if err != nil {
			return database.WrapError(op, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return database.WrapError(op, err)
	}
	return nil
}

func (r *TeamRepository) GetTeamByName(ctx context.Context, name string) (*models.Team, error) {
	const op = "GetTeamByName"

	exists, err := r.TeamExists(ctx, name)
	if err != nil {
		return nil, database.WrapError(op, err)
	}
	if !exists {
		return nil, database.WrapError(op, database.ErrTeamNotFound)
	}

	query := `
		SELECT user_id, username, is_active 
		FROM users 
		WHERE team_name = ? 
		ORDER BY username
	`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, database.WrapError(op, err)
	}
	defer rows.Close()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		err := rows.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return nil, database.WrapError(op, err)
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, database.WrapError(op, err)
	}

	team := &models.Team{
		Name:    name,
		Members: members,
	}

	return team, nil
}

func (r *TeamRepository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	const op = "TeamExists"

	query := `SELECT 1 FROM teams WHERE name = ?`
	row := r.db.QueryRowContext(ctx, query, teamName)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, database.WrapError(op, err)
	}

	return true, nil
}
