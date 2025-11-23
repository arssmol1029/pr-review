package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pr-review/internal/config"
	"pr-review/internal/errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func New(ctx context.Context, cfg config.DatabaseConfig) (*PostgresRepository, error) {
	const op = "PostgresRepository.New"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.WrapError(op, err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.WrapError(op, err)
	}

	r := &PostgresRepository{db: db}

	if err := r.runMigrations(cfg.MigrationsPath); err != nil {
		return nil, errors.WrapError(op, err)
	}

	log.Println("PostgreSQL repository initialized successfully")
	return r, nil
}

func (r *PostgresRepository) runMigrations(migrationsPath string) error {
	const op = "PostgresRepository.runMigrations"

	driver, err := postgres.WithInstance(r.db, &postgres.Config{})
	if err != nil {
		return errors.WrapError(op, err)
	}

	sourceURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("%s: failed to create migration instance: %w", op, err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.WrapError(op, err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func (r *PostgresRepository) Close() error {
	const op = "PostgresRepository.Close"

	if r.db != nil {
		err := r.db.Close()
		if err != nil {
			return errors.WrapError(op, err)
		}
	}
	return nil
}

func (r *PostgresRepository) Ping(ctx context.Context) error {
	const op = "PostgresRepository.Ping"

	if r.db == nil {
		return fmt.Errorf("%s: database not initialized", op)
	}

	err := r.db.PingContext(ctx)
	if err != nil {
		return errors.WrapError(op, err)
	}
	return nil
}
