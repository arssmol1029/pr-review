package service

import (
	"context"
	"log/slog"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/server/handlers"
)

type UserRepository interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
}

type userService struct {
	logger *slog.Logger
	repo   UserRepository
}

func NewuserService(
	logger *slog.Logger,
	repo UserRepository,
) handlers.UserService {
	return &userService{
		logger: logger,
		repo:   repo,
	}
}

func (s *userService) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	const op = "userService.SetUserActive"

	err := s.repo.SetUserActive(ctx, userID, isActive)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("failed to set user active", "error", err, "userID", userID, "isActive", isActive)
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("failed to get updated user", "error", err, "userID", userID)
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserReviewPRs(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	const op = "userService.GetUserReviewPRs"

	prs, err := s.repo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("failed to get user review PRs", "error", err, "userID", userID)
		return nil, err
	}

	return prs, nil
}
