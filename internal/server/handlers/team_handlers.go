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
	"github.com/go-playground/validator/v10"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type TeamHandler struct {
	logger  *slog.Logger
	service TeamService
}

func NewTeamHandler(logger *slog.Logger, s TeamService) *TeamHandler {
	return &TeamHandler{
		logger:  logger,
		service: s,
	}
}

// POST /team/add
func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	const op = "TeamHandlers.Add"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	type MemberItem struct {
		UserID   string `json:"user_id" validate:"required"`
		Username string `json:"username" validate:"required"`
		IsActive bool   `json:"is_active"`
	}

	var req struct {
		Name    string       `json:"team_name" validate:"required"`
		Members []MemberItem `json:"members" validate:"dive"`
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("Failed to decode request body", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ERROR("INVALID_REQUEST", "wrong request format"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error("Request validation failed", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ERROR("VALIDATION_ERROR", "wrong request format"))
		return
	}

	members := make([]models.TeamMember, 0)
	for _, m := range req.Members {
		member := models.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
		members = append(members, member)
	}
	team := &models.Team{
		Name:    req.Name,
		Members: members,
	}

	createdTeam, err := h.service.CreateTeam(r.Context(), team)
	if errors.Is(err, serviceErrors.ErrTeamExists) {
		log.Error("Team already exists", "error", err, "team_name", team.Name)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.TEAM_EXISTS())
		return
	}
	if errors.Is(err, serviceErrors.ErrUserExists) {
		log.Error("User already exists", "error", err, "team_name", team.Name)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.USER_EXISTS())
		return
	}
	if err != nil {
		log.Error("Failed to create team", "error", err, "team_name", team.Name)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to create team"))
		return
	}

	type TeamItem struct {
		Name    string       `json:"team_name" validate:"required"`
		Members []MemberItem `json:"members,omitempty" validate:"dive"`
	}

	res := struct {
		Team TeamItem `json:"team" validate:"required"`
	}{
		Team: TeamItem{
			Name: createdTeam.Name,
		},
	}

	for _, m := range createdTeam.Members {
		member := MemberItem{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
		res.Team.Members = append(res.Team.Members, member)
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}

// GET /team/get
func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	const op = "TeamHandlers.Get"

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

	team, err := h.service.GetTeam(r.Context(), teamName)
	if errors.Is(err, serviceErrors.ErrTeamNotFound) {
		log.Error("Team not found", "error", err, "team_name", teamName)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("team not found"))
		return
	}
	if err != nil {
		log.Error("Failed to get team", "error", err, "team_name", teamName)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to get team"))
		return
	}

	type MemberItem struct {
		UserID   string `json:"user_id" validate:"required"`
		Username string `json:"username" validate:"required"`
		IsActive bool   `json:"is_active"`
	}

	res := struct {
		Name    string       `json:"team_name" validate:"required"`
		Members []MemberItem `json:"members,omitempty" validate:"dive"`
	}{
		Name: team.Name,
	}

	for _, m := range team.Members {
		member := MemberItem{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
		res.Members = append(res.Members, member)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}
