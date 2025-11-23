package service

import (
	"context"
	"log/slog"

	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/server/handlers"
)

type StatsRepository interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetPRsByReviewer(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
	GetPRsCntByAuthor(ctx context.Context, userID string) (int, error)
	GetTeamByName(ctx context.Context, name string) (*models.Team, error)
	GetPRsCntByTeam(ctx context.Context, teamName string) (int, error)
	GetAvgReviewersPerPR(ctx context.Context, teamName string) (float64, error)
	GetTotalStats(ctx context.Context) (*models.TotalStats, error)
}

type statsService struct {
	logger *slog.Logger
	repo   StatsRepository
}

func NewStatsService(
	logger *slog.Logger,
	repo StatsRepository,
) handlers.StatsService {
	return &statsService{
		logger: logger,
		repo:   repo,
	}
}

func (s *statsService) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	const op = "statsService.GetUserStats"

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get user", "error", err, "userID", userID)
		return nil, err
	}

	stats := &models.UserStats{}
	stats.UserID = user.UserID
	stats.Username = user.Username
	stats.TeamName = user.TeamName

	prs, err := s.repo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get user review PRs", "error", err, "userID", userID)
		return nil, err
	}

	for _, pr := range prs {
		switch pr.Status {
		case "OPEN":
			stats.OpenReviews++
		case "MERGED":
			stats.MergedReviews++
		}
	}

	count, err := s.repo.GetPRsCntByAuthor(ctx, userID)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get user authored PRs count", "error", err, "userID", userID)
		return nil, err
	}
	stats.CreatedPRs = count

	return stats, nil
}

func (s *statsService) GetTeamStats(ctx context.Context, teamName string) (*models.TeamStats, error) {
	const op = "statsService.GetTeamStats"

	team, err := s.repo.GetTeamByName(ctx, teamName)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get team", "error", err, "teamName", teamName)
		return nil, err
	}

	stats := &models.TeamStats{}
	stats.TeamName = team.Name

	for _, user := range team.Members {
		stats.MemberCount++
		if user.IsActive {
			stats.ActiveMembers++
		}
	}

	prsCount, err := s.repo.GetPRsCntByTeam(ctx, teamName)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get team created PRs count", "error", err, "teamName", teamName)
		return nil, err
	}
	stats.CreatedPRs = prsCount

	avgReviewersCount, err := s.repo.GetAvgReviewersPerPR(ctx, teamName)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get team average reviewers count", "error", err, "teamName", teamName)
		return nil, err
	}
	stats.AvgReviewersPerPR = avgReviewersCount

	return stats, nil
}

func (s *statsService) GetTotalStats(ctx context.Context) (*models.TotalStats, error) {
	const op = "statsService.GetTotalStats"

	stats, err := s.repo.GetTotalStats(ctx)
	if err != nil {
		err = errors.WrapError(op, err)
		s.logger.Error("Failed to get stats", "error", err)
		return nil, err
	}

	return stats, nil
}
