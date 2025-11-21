package database

import (
	"context"
	"pr-review/internal/database/models"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	DeleteUser(ctx context.Context, id string) error
	GetUsersByTeam(ctx context.Context, teamName string) ([]*models.User, error)
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]*models.User, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, name string) (*models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) error
	DeleteTeam(ctx context.Context, name string) error
}

type PRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequest) error
	GetPRByID(ctx context.Context, id string) (*models.PullRequest, error)
	UpdatePR(ctx context.Context, pr *models.PullRequest) error
	MergePR(ctx context.Context, prID string, mergedAt time.Time) error
	AssignReviewer(ctx context.Context, prID, userID string) error
	RemoveReviewer(ctx context.Context, prID, userID string) error
	GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequest, error)
	GetPRsByAuthor(ctx context.Context, authorID string) ([]*models.PullRequest, error)
	PRExists(ctx context.Context, prID string) (bool, error)
}

type RepositoryCollection struct {
	Users UserRepository
	Teams TeamRepository
	PRs   PRRepository
}
