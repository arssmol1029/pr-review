package service

import (
	"log/slog"
	"pr-review/internal/database"
)

type PRService interface{}

type prService struct {
	logger   *slog.Logger
	userRepo *database.UserRepository
	teamRepo *database.TeamRepository
	prRepo   *database.PRRepository
}

func NewPRService(
	logger *slog.Logger,
	userRepo *database.UserRepository,
	teamRepo *database.TeamRepository,
	prRepo *database.PRRepository,
) PRService {
	return &prService{
		logger:   logger,
		userRepo: userRepo,
		teamRepo: teamRepo,
		prRepo:   prRepo,
	}
}
