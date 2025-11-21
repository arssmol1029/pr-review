package handlers

import (
	"log/slog"
	"net/http"
	"pr-review/internal/storage"
)

type PRHandlers struct {
	logger   *slog.Logger
	userRepo *storage.UserRepository
	teamRepo *storage.TeamRepository
	prRepo   *storage.PRRepository
}

func NewPR(
	logger *slog.Logger,
	userRepo *storage.UserRepository,
	teamRepo *storage.TeamRepository,
	prRepo *storage.PRRepository,
) *PRHandlers {
	return &PRHandlers{
		logger:   logger,
		userRepo: userRepo,
		teamRepo: teamRepo,
		prRepo:   prRepo,
	}
}

// POST /pullRequest/create
func (*PRHandlers) Create(w http.ResponseWriter, r *http.Request) {

}

// POST /pullRequest/merge
func (*PRHandlers) Merge(w http.ResponseWriter, r *http.Request) {

}

// POST /pullRequest/reassign
func (*PRHandlers) Reassign(w http.ResponseWriter, r *http.Request) {

}
