package handlers

import (
	"context"
	"net/http"
	"pr-review/internal/models"
)

type UserService interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}

// POST /users/setIsActive
func (*UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {

}

// GET /users/getReview
func (*UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {

}
