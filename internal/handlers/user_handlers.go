package handlers

import (
	"log/slog"
	"net/http"
	"pr-review/internal/storage"
)

type UserHandlers struct {
	logger   *slog.Logger
	userRepo *storage.UserRepository
}

func NewUser(
	logger *slog.Logger,
	userRepo *storage.UserRepository,
) *UserHandlers {
	return &UserHandlers{
		logger:   logger,
		userRepo: userRepo,
	}
}

// POST /users/setIsActive
func (*UserHandlers) SetIsActive(w http.ResponseWriter, r *http.Request) {

}

// GET /users/getReview
func (*UserHandlers) GetReview(w http.ResponseWriter, r *http.Request) {

}
