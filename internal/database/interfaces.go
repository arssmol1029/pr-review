package database

import (
	"context"
	"pr-review/internal/database/models"
	"time"
)

type Database interface {
	Init(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequest, error)
	UserExists(ctx context.Context, userID string) (bool, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, name string) (*models.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}

type PRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequest) error
	GetPRByID(ctx context.Context, id string) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string, mergedAt time.Time) error
	ReassignReviewer(ctx context.Context, prID, userID string) error
	PRExists(ctx context.Context, prID string) (bool, error)
}

type Repositories struct {
	Users UserRepository
	Teams TeamRepository
	PRs   PRRepository
}
