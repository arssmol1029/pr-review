package service

import (
	"context"
	"log/slog"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"time"
)

type PRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequestShort) error
	GetPRByID(ctx context.Context, id string) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string, mergedAt time.Time) error
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*string, error)
}

type PRService struct {
	logger *slog.Logger
	prRepo PRRepository
}

func NewPRService(
	logger *slog.Logger,
	prRepo PRRepository,
) *PRService {
	return &PRService{
		logger: logger,
		prRepo: prRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, pr *models.PullRequestShort) (*models.PullRequest, error) {
	const op = "PRService.CreatePR"

	err := s.prRepo.CreatePR(ctx, pr)
	if err != nil {
		s.logger.Error("failed to create PR", "op", op, "error", err, "prID", pr.ID)
		return nil, errors.WrapError(op, err)
	}

	createdPR, err := s.prRepo.GetPRByID(ctx, pr.ID)
	if err != nil {
		s.logger.Error("failed to get created PR", "op", op, "error", err, "prID", pr.ID)
		return nil, errors.WrapError(op, err)
	}

	return createdPR, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	const op = "PRService.MergePR"

	err := s.prRepo.MergePR(ctx, prID, time.Now())
	if err != nil {
		s.logger.Error("failed to merge PR", "op", op, "error", err, "prID", prID)
		return nil, errors.WrapError(op, err)
	}

	mergedPR, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		s.logger.Error("failed to get merged PR", "op", op, "error", err, "prID", prID)
		return nil, errors.WrapError(op, err)
	}

	return mergedPR, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, *string, error) {
	const op = "PRService.ReassignReviewer"

	newUserID, err := s.prRepo.ReassignReviewer(ctx, prID, oldUserID)
	if err != nil {
		s.logger.Error("failed to reassign reviewer", "op", op, "error", err, "prID", prID, "oldUserID", oldUserID)
		return nil, nil, errors.WrapError(op, err)
	}

	updatedPR, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		s.logger.Error("failed to get updated PR", "op", op, "error", err, "prID", prID)
		return nil, nil, errors.WrapError(op, err)
	}

	return updatedPR, newUserID, nil
}
