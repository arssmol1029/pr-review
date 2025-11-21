package service

import (
	"log/slog"
	"pr-review/internal/database"
)

type TeamService interface{}

type teamService struct {
	logger   *slog.Logger
	userRepo *database.UserRepository
	teamRepo *database.TeamRepository
}

func NewTeamService(
	logger *slog.Logger,
	userRepo *database.UserRepository,
	teamRepo *database.TeamRepository,
) TeamService {
	return &teamService{
		logger:   logger,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}
