package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	serviceErrors "pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/server/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type StatsService interface {
	GetTeamStats(ctx context.Context, userID string) (*models.TeamStats, error)
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
	GetTotalStats(ctx context.Context) (*models.TotalStats, error)
}

type StatsHandler struct {
	logger  *slog.Logger
	service StatsService
}

func NewStatsHandler(logger *slog.Logger, s StatsService) *StatsHandler {
	return &StatsHandler{
		logger:  logger,
		service: s,
	}
}

// GET stats/user
func (h *StatsHandler) User(w http.ResponseWriter, r *http.Request) {
	const op = "StatsHandlers.Stats"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Error("Missing userID parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ERROR("INVALID_REQUEST", "user_id query parameter is required"))
		return
	}

	stats, err := h.service.GetUserStats(r.Context(), userID)
	if errors.Is(err, serviceErrors.ErrUserNotFound) {
		log.Error("User not found", "error", err, "userID", userID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("user not found"))
		return
	}
	if err != nil {
		log.Error("Failed to get user stats", "error", err, "userID", userID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to get user stats"))
		return
	}

	type StatsItem struct {
		UserID        string `json:"user_id" validate:"required"`
		Username      string `json:"username" validate:"required"`
		TeamName      string `json:"team_name" validate:"required"`
		OpenReviews   int    `json:"open_assignments" validate:"required"`
		MergedReviews int    `json:"merged_assignments" validate:"required"`
		CreatedPRs    int    `json:"created_prs" validate:"required"`
	}

	res := struct {
		Stats StatsItem `json:"user_stats" validate:"required"`
	}{
		Stats: StatsItem{
			UserID:        stats.UserID,
			Username:      stats.Username,
			TeamName:      stats.TeamName,
			OpenReviews:   stats.OpenReviews,
			MergedReviews: stats.MergedReviews,
			CreatedPRs:    stats.CreatedPRs,
		},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GET stats/team
func (h *StatsHandler) Team(w http.ResponseWriter, r *http.Request) {
	const op = "StatsHandlers.Stats"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		log.Error("team_name query parameter is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ERROR("INVALID_REQUEST", "team_name query parameter is required"))
		return
	}

	stats, err := h.service.GetTeamStats(r.Context(), teamName)
	if errors.Is(err, serviceErrors.ErrTeamNotFound) {
		log.Error("Team not found", "error", err, "team_name", teamName)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("team not found"))
		return
	}
	if err != nil {
		log.Error("Failed to get team stats", "error", err, "team_name", teamName)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to get team stats"))
		return
	}

	type StatsItem struct {
		TeamName          string  `json:"team_name" validate:"required"`
		MemberCount       int     `json:"member_count" validate:"required"`
		ActiveMembers     int     `json:"active_members" validate:"required"`
		CreatedPRs        int     `json:"created_prs" validate:"required"`
		AvgReviewersPerPR float64 `json:"avg_reviewers_per_pr" validate:"required"`
	}

	res := struct {
		Stats StatsItem `json:"team_stats" validate:"required"`
	}{
		Stats: StatsItem{
			TeamName:          stats.TeamName,
			MemberCount:       stats.MemberCount,
			ActiveMembers:     stats.ActiveMembers,
			CreatedPRs:        stats.CreatedPRs,
			AvgReviewersPerPR: stats.AvgReviewersPerPR,
		},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GET stats/total
func (h *StatsHandler) Total(w http.ResponseWriter, r *http.Request) {
	const op = "StatsHandlers.Stats"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	stats, err := h.service.GetTotalStats(r.Context())
	if err != nil {
		log.Error("Failed to get stats", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to get stats"))
		return
	}

	type StatsItem struct {
		TotalTeams        int     `json:"total_teams" validate:"required"`
		TotalUsers        int     `json:"total_users" validate:"required"`
		ActiveUsers       int     `json:"active_users" validate:"required"`
		TotalPRs          int     `json:"total_prs" validate:"required"`
		OpenPRs           int     `json:"open_prs" validate:"required"`
		MergedPRs         int     `json:"merged_prs" validate:"required"`
		AvgReviewersPerPR float64 `json:"avg_reviewers_per_pr" validate:"required"`
	}

	res := struct {
		Stats StatsItem `json:"total_stats" validate:"required"`
	}{
		Stats: StatsItem{
			TotalTeams:        stats.TotalTeams,
			TotalUsers:        stats.TotalUsers,
			ActiveUsers:       stats.ActiveUsers,
			TotalPRs:          stats.TotalPRs,
			OpenPRs:           stats.OpenPRs,
			MergedPRs:         stats.MergedPRs,
			AvgReviewersPerPR: stats.AvgReviewersPerPR,
		},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}
