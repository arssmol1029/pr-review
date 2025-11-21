package service

import (
	"log/slog"
	"pr-review/internal/database"
)

type UserService interface{}

type userService struct {
	logger   *slog.Logger
	userRepo *database.UserRepository
}

func NewUserService(
	logger *slog.Logger,
	userRepo *database.UserRepository,
) UserService {
	return &userService{
		logger:   logger,
		userRepo: userRepo,
	}
}
