package service

import (
	"context"
	"log/slog"
	"pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/server/handlers"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, name string) (*models.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
	GetPRsCntByTeam(ctx context.Context, teamName string) (int, error)
	GetAvgReviewersPerPR(ctx context.Context, teamName string) (float64, error)
}

type teamService struct {
	logger *slog.Logger
	repo   TeamRepository
}

func NewTeamService(
	logger *slog.Logger,
	repo TeamRepository,
) handlers.TeamService {
	return &teamService{
		logger: logger,
		repo:   repo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	const op = "teamService.CreateTeam"

	err := s.repo.CreateTeam(ctx, team)
	if err != nil {
		s.logger.Error("failed to create team", "op", op, "error", err, "teamName", team.Name)
		return nil, errors.WrapError(op, err)
	}

	createdTeam, err := s.repo.GetTeamByName(ctx, team.Name)
	if err != nil {
		s.logger.Error("failed to get created team", "op", op, "error", err, "teamName", team.Name)
		return nil, errors.WrapError(op, err)
	}

	return createdTeam, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	const op = "teamService.GetTeam"

	team, err := s.repo.GetTeamByName(ctx, teamName)
	if err != nil {
		s.logger.Error("failed to get team", "op", op, "error", err, "teamName", teamName)
		return nil, errors.WrapError(op, err)
	}

	return team, nil
}

func (s *teamService) GetTeamStats(ctx context.Context, teamName string) (*models.TeamStats, error) {
	const op = "userService.GetTeamStats"

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
