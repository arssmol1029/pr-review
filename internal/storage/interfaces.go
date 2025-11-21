package storage

import (
	"context"
	"pr-review/internal/types"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *types.User) error
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetUserByUsername(ctx context.Context, email string) (*types.User, error)
	UpdateUser(ctx context.Context, user *types.User) error
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	DeleteUser(ctx context.Context, id string) error
	GetUsersByTeam(ctx context.Context, teamName string) ([]*types.User, error)
	GetActiveUsersByTeam(ctx context.Context, teamName string) ([]*types.User, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *types.Team) error
	GetTeamByName(ctx context.Context, name string) (*types.Team, error)
	UpdateTeam(ctx context.Context, team *types.Team) error
	DeleteTeam(ctx context.Context, name string) error
}

type PRRepository interface {
	CreatePR(ctx context.Context, pr *types.PullRequest) error
	GetPRByID(ctx context.Context, id string) (*types.PullRequest, error)
	UpdatePR(ctx context.Context, pr *types.PullRequest) error
	MergePR(ctx context.Context, prID string, mergedAt time.Time) error
	AssignReviewer(ctx context.Context, prID, userID string) error
	RemoveReviewer(ctx context.Context, prID, userID string) error
	GetPRsByReviewer(ctx context.Context, userID string) ([]*types.PullRequest, error)
	GetPRsByAuthor(ctx context.Context, authorID string) ([]*types.PullRequest, error)
	PRExists(ctx context.Context, prID string) (bool, error)
}

type RepositoryCollection struct {
	Users UserRepository
	Teams TeamRepository
	PRs   PRRepository
}
