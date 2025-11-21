package handlers

import (
	"log/slog"
	"net/http"
	"pr-review/internal/storage"
)

type TeamHandlers struct {
	logger   *slog.Logger
	userRepo *storage.UserRepository
	teamRepo *storage.TeamRepository
}

func NewTeam(
	logger *slog.Logger,
	userRepo *storage.UserRepository,
	teamRepo *storage.TeamRepository,
) *TeamHandlers {
	return &TeamHandlers{
		logger:   logger,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// POST /team/add
func (*TeamHandlers) Add(w http.ResponseWriter, r *http.Request) {

}

// GET /team/get
func (*TeamHandlers) Get(w http.ResponseWriter, r *http.Request) {

}
