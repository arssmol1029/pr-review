package service

import (
	"context"
	"log/slog"
	"pr-review/internal/errors"
	"pr-review/internal/models"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, name string) (*models.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)
}

type TeamService struct {
	logger   *slog.Logger
	teamRepo TeamRepository
}

func NewTeamService(
	logger *slog.Logger,
	teamRepo TeamRepository,
) *TeamService {
	return &TeamService{
		logger:   logger,
		teamRepo: teamRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	const op = "TeamService.CreateTeam"

	err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		s.logger.Error("failed to create team", "op", op, "error", err, "teamName", team.Name)
		return nil, errors.WrapError(op, err)
	}

	createdTeam, err := s.teamRepo.GetTeamByName(ctx, team.Name)
	if err != nil {
		s.logger.Error("failed to get created team", "op", op, "error", err, "teamName", team.Name)
		return nil, errors.WrapError(op, err)
	}

	return createdTeam, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	const op = "TeamService.GetTeam"

	team, err := s.teamRepo.GetTeamByName(ctx, teamName)
	if err != nil {
		s.logger.Error("failed to get team", "op", op, "error", err, "teamName", teamName)
		return nil, errors.WrapError(op, err)
	}

	return team, nil
}
