package service

import (
	"context"
	"log/slog"
	"pr-review/internal/database/models"
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
