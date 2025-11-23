package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *sql.DB
}

func New(ctx context.Context, storagePath string) (*SQLiteRepository, error) {
	const op = "SQLiteDatabase.Init"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	r := &SQLiteRepository{db: db}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS teams (
			name TEXT PRIMARY KEY
		)`,

		`CREATE TABLE IF NOT EXISTS users (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			team_name TEXT NOT NULL,
			FOREIGN KEY (team_name) REFERENCES teams(name) ON DELETE CASCADE
		)`,

		`CREATE TABLE IF NOT EXISTS pull_requests (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			author_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'OPEN',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			merged_at DATETIME DEFAULT NULL,
			FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,

		`CREATE TABLE IF NOT EXISTS pr_reviewers (
			pr_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (pr_id, user_id),
			FOREIGN KEY (pr_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)`,

		`CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_prs_author ON pull_requests(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pr_reviewers(user_id)`,
	}

	for _, query := range queries {
		if _, err := r.db.ExecContext(ctx, query); err != nil {
			return nil, fmt.Errorf("%s: failed to execute query %s: %w", op, query, err)
		}
	}

	return r, nil
}

func (r *SQLiteRepository) Close() error {
	const op = "SQLiteDatabase.Close"

	if r.db != nil {
		err := r.db.Close()
		if err != nil {
			return fmt.Errorf("%s: failed to close database: %w", op, err)
		}
	}
	return nil
}

func (r *SQLiteRepository) Ping(ctx context.Context) error {
	const op = "SQLiteDatabase.Ping"

	if r.db == nil {
		return fmt.Errorf("%s: database not initialized", op)
	}

	err := r.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to ping database: %w", op, err)
	}
	return nil
}
