package service

import (
	"context"
	"log/slog"
	"pr-review/internal/database/models"
	"time"
)

type PRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequestShort) error
	GetPRByID(ctx context.Context, id string) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string, mergedAt time.Time) error
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*string, error)
	PRExists(ctx context.Context, prID string) (bool, error)
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
