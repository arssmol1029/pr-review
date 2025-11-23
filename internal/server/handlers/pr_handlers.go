package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	serviceErrors "pr-review/internal/errors"
	"pr-review/internal/models"
	"pr-review/internal/server/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type PRService interface {
	CreatePR(ctx context.Context, pr *models.PullRequestShort) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, *string, error)
}

type PRHandler struct {
	logger  *slog.Logger
	service PRService
}

func NewPRHandler(logger *slog.Logger, s PRService) *PRHandler {
	return &PRHandler{
		logger:  logger,
		service: s,
	}
}

// POST /pullRequest/create
func (h *PRHandler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "PRHandler.Create"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req struct {
		ID       string `json:"pull_request_id" validate:"required"`
		Name     string `json:"pull_request_name" validate:"required"`
		AuthorID string `json:"author_id" validate:"required"`
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

	prShort := &models.PullRequestShort{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
		Status:   "OPEN",
	}
	pr, err := h.service.CreatePR(r.Context(), prShort)
	if errors.Is(err, serviceErrors.ErrUserNotFound) {
		log.Error("Author not found", "error", err, "author_id", req.AuthorID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("author not found"))
		return
	}
	if errors.Is(err, serviceErrors.ErrPRExists) {
		log.Error("PR already exists", "error", err, "prID", req.ID)
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.PR_EXISTS())
		return
	}
	if err != nil {
		log.Error("Failed to create PR", "error", err, "prID", req.ID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to create pull request"))
		return
	}

	type PRItem struct {
		ID                string   `json:"pull_request_id" validate:"required"`
		Name              string   `json:"pull_request_name" validate:"required"`
		AuthorID          string   `json:"author_id" validate:"required"`
		Status            string   `json:"status" validate:"required"`
		AssignedReviewers []string `json:"assigned_reviewers"`
	}

	res := struct {
		PullRequest PRItem `json:"pr" validate:"required"`
	}{
		PullRequest: PRItem{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
		},
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}

// POST /pullRequest/merge
func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	const op = "PRHandler.Merge"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req struct {
		ID string `json:"pull_request_id" validate:"required"`
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

	pr, err := h.service.MergePR(r.Context(), req.ID)
	if errors.Is(err, serviceErrors.ErrPRNotFound) {
		log.Error("PR not found", "error", err, "prID", req.ID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("pull request not found"))
		return
	}
	if err != nil {
		log.Error("Failed to merge PR", "error", err, "prID", req.ID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to merge pull request"))
		return
	}

	type PRItem struct {
		MetgedAt          *time.Time `json:"merged_at" validate:"required"`
		ID                string     `json:"pull_request_id" validate:"required"`
		Name              string     `json:"pull_request_name" validate:"required"`
		AuthorID          string     `json:"author_id" validate:"required"`
		Status            string     `json:"status" validate:"required"`
		AssignedReviewers []string   `json:"assigned_reviewers" validate:"required"`
	}

	res := struct {
		PullRequest PRItem `json:"pr" validate:"required"`
	}{
		PullRequest: PRItem{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
			MetgedAt:          pr.MergedAt,
		},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}

// POST /pullRequest/reassign
func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	const op = "PRHandler.Reassign"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req struct {
		PullRequestID string `json:"pull_request_id" validate:"required"`
		OldUserID     string `json:"old_user_id" validate:"required"`
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

	pr, newUserID, err := h.service.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if errors.Is(err, serviceErrors.ErrPRNotFound) {
		log.Error("PR not found", "error", err, "prID", req.PullRequestID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("pull request not found"))
		return
	}
	if errors.Is(err, serviceErrors.ErrUserNotFound) {
		log.Error("Reviewer not found", "error", err, "old_user_id", req.OldUserID, "prID", req.PullRequestID)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.NOT_FOUND("reviewer not found"))
		return
	}
	if errors.Is(err, serviceErrors.ErrPRMerged) {
		log.Error("PR already merged", "error", err, "prID", req.PullRequestID)
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.PR_MERGED())
		return
	}
	if errors.Is(err, serviceErrors.ErrNotAssigned) {
		log.Error("Reviewer not assigned to this PR", "error", err, "old_user_id", req.OldUserID, "prID", req.PullRequestID)
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.NOT_ASSIGNED())
		return
	}
	if errors.Is(err, serviceErrors.ErrNoCandidate) {
		log.Error("No active replacement candidate in team", "error", err, "old_user_id", req.OldUserID, "prID", req.PullRequestID)
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.NO_CANDIDATE())
		return
	}
	if err != nil {
		log.Error("Failed to reassign reviewer", "error", err, "prID", req.PullRequestID, "old_user_id", req.OldUserID)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ERROR("INTERNAL_ERROR", "failed to reassign reviewer"))
		return
	}

	type PRItem struct {
		ID                string   `json:"pull_request_id" validate:"required"`
		Name              string   `json:"pull_request_name" validate:"required"`
		AuthorID          string   `json:"author_id" validate:"required"`
		Status            string   `json:"status" validate:"required"`
		AssignedReviewers []string `json:"assigned_reviewers" validate:"required"`
	}

	res := struct {
		NewUserID   string `json:"replaced_by" validate:"required"`
		PullRequest PRItem `json:"pr" validate:"required"`
	}{
		PullRequest: PRItem{
			ID:                pr.ID,
			Name:              pr.Name,
			AuthorID:          pr.AuthorID,
			Status:            pr.Status,
			AssignedReviewers: pr.AssignedReviewers,
		},
		NewUserID: *newUserID,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, res)
}
