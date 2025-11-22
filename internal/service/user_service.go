package service

import (
	"context"
	"log/slog"
	"pr-review/internal/errors"
	"pr-review/internal/models"
)

type UserRepository interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
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

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	const op = "UserService.SetUserActive"

	err := s.userRepo.SetUserActive(ctx, userID, isActive)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("failed to set user active", "error", err, "userID", userID, "isActive", isActive)
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("failed to get updated user", "error", err, "userID", userID)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserReviewPRs(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	const op = "UserService.GetUserReviewPRs"

	prs, err := s.userRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("failed to get user review PRs", "error", err, "userID", userID)
		return nil, err
	}

	return prs, nil
}
