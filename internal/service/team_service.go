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
		s.logger.Error("Failed to create team", "op", op, "error", err, "teamName", team.Name)
		return nil, errors.WrapError(op, err)
	}

	createdTeam, err := s.repo.GetTeamByName(ctx, team.Name)
	if err != nil {
		s.logger.Error("Failed to get created team", "op", op, "error", err, "teamName", team.Name)
		return nil, errors.WrapError(op, err)
	}

	return createdTeam, nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	const op = "teamService.GetTeam"

	team, err := s.repo.GetTeamByName(ctx, teamName)
	if err != nil {
		s.logger.Error("Failed to get team", "op", op, "error", err, "teamName", teamName)
		return nil, errors.WrapError(op, err)
	}

	return team, nil
}
