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

type UserService interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
}

type UserHandler struct {
	logger  *slog.Logger
	service UserService
}

func NewUserHandler(logger *slog.Logger, s UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: s,
	}
}

// POST /users/setIsActive
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandlers.SetIsActive"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req struct {
		UserID   string `json:"user_id" validate:"required"`
		IsActive bool   `json:"is_active"`
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

	user, err := h.service.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if errors.Is(err, serviceErrors.ErrUserNotFound) {
		log.Error("User not found", "error", err, "user_id", req.UserID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("user not found"))
		return
	}
	if err != nil {
		log.Error("Failed to set user active status", "error", err, "user_id", req.UserID, "is_active", req.IsActive)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to set user active status"))
		return
	}

	type UserItem struct {
		UserID   string `json:"user_id" validate:"required"`
		Username string `json:"username" validate:"required"`
		TeamName string `json:"team_name" validate:"required"`
		IsActive bool   `json:"is_active"`
	}

	res := struct {
		User UserItem `json:"user" validate:"required"`
	}{
		User: UserItem{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}

	// if err := validate.Struct(res); err != nil {
	// 	log.Error("Response validation failed", "error", err)
	// 	render.Status(r, http.StatusInternalServerError)
	// 	render.JSON(w, r, response.ERROR("VALIDATION_ERROR", "wrong response format"))
	// 	return
	// }

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// GET /users/getReview
func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandlers.GetReview"

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

	prs, err := h.service.GetUserReviewPRs(r.Context(), userID)
	if errors.Is(err, serviceErrors.ErrUserNotFound) {
		log.Error("User not found", "error", err, "userID", userID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("user not found"))
		return
	}
	if err != nil {
		log.Error("Failed to get user review PRs", "error", err, "userID", userID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to get user review PRs"))
		return
	}

	type PRItem struct {
		ID       string `json:"pull_request_id" validate:"required"`
		PRName   string `json:"pull_request_name" validate:"required"`
		AuthorID string `json:"author_id" validate:"required"`
		Status   string `json:"status" validate:"required"`
	}

	var res struct {
		UserID       string   `json:"user_id" validate:"required"`
		PullRequests []PRItem `json:"pull_requests" validate:"required,dive"`
	}

	res.UserID = userID
	for _, pr := range prs {
		prRes := PRItem{
			ID:       pr.ID,
			PRName:   pr.Name,
			AuthorID: pr.AuthorID,
			Status:   pr.Status,
		}
		res.PullRequests = append(res.PullRequests, prRes)
	}

	// validate := validator.New()
	// if err := validate.Struct(res); err != nil {
	// 	log.Error("Response validation failed", "error", err)
	// 	render.Status(r, http.StatusInternalServerError)
	// 	render.JSON(w, r, response.ERROR("VALIDATION_ERROR", "wrong response format"))
	// 	return
	// }

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}
