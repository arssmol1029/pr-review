package service

import (
	"context"
	"log/slog"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/server/handlers"
	"time"
)

type PRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequestShort) error
	GetPRByID(ctx context.Context, id string) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string, mergedAt time.Time) error
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*string, error)
}

type prService struct {
	logger *slog.Logger
	repo   PRRepository
}

func NewprService(
	logger *slog.Logger,
	repo PRRepository,
) handlers.PRService {
	return &prService{
		logger: logger,
		repo:   repo,
	}
}

func (s *prService) CreatePR(ctx context.Context, pr *models.PullRequestShort) (*models.PullRequest, error) {
	const op = "prService.CreatePR"

	err := s.repo.CreatePR(ctx, pr)
	if err != nil {
		s.logger.Error("failed to create PR", "op", op, "error", err, "prID", pr.ID)
		return nil, errors.WrapError(op, err)
	}

	createdPR, err := s.repo.GetPRByID(ctx, pr.ID)
	if err != nil {
		s.logger.Error("failed to get created PR", "op", op, "error", err, "prID", pr.ID)
		return nil, errors.WrapError(op, err)
	}

	return createdPR, nil
}

func (s *prService) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	const op = "prService.MergePR"

	err := s.repo.MergePR(ctx, prID, time.Now())
	if err != nil {
		s.logger.Error("failed to merge PR", "op", op, "error", err, "prID", prID)
		return nil, errors.WrapError(op, err)
	}

	mergedPR, err := s.repo.GetPRByID(ctx, prID)
	if err != nil {
		s.logger.Error("failed to get merged PR", "op", op, "error", err, "prID", prID)
		return nil, errors.WrapError(op, err)
	}

	return mergedPR, nil
}

func (s *prService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, *string, error) {
	const op = "prService.ReassignReviewer"

	newUserID, err := s.repo.ReassignReviewer(ctx, prID, oldUserID)
	if err != nil {
		s.logger.Error("failed to reassign reviewer", "op", op, "error", err, "prID", prID, "oldUserID", oldUserID)
		return nil, nil, errors.WrapError(op, err)
	}

	updatedPR, err := s.repo.GetPRByID(ctx, prID)
	if err != nil {
		s.logger.Error("failed to get updated PR", "op", op, "error", err, "prID", prID)
		return nil, nil, errors.WrapError(op, err)
	}

	return updatedPR, newUserID, nil
}
