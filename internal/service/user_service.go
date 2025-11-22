package service

import (
	"context"
	"log/slog"
	"pr-review/internal/database/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
	UserExists(ctx context.Context, userID string) (bool, error)
}

type UserService struct {
	logger   *slog.Logger
	userRepo UserRepository
}

func NewUserService(
	logger *slog.Logger,
	userRepo UserRepository,
) *UserService {
	return &UserService{
		logger:   logger,
		userRepo: userRepo,
	}
}
