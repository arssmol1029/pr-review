package handlers

import (
	"context"
	"net/http"
	"pr-review/internal/models"
)

type PRService interface {
	CreatePR(ctx context.Context, pr *models.PullRequestShort) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*models.PullRequest, *string, error)
}

type PRHandler struct {
	service PRService
}

func NewPRHandler(s PRService) *PRHandler {
	return &PRHandler{service: s}
}

// POST /pullRequest/create
func (*PRHandler) Create(w http.ResponseWriter, r *http.Request) {

}

// POST /pullRequest/merge
func (*PRHandler) Merge(w http.ResponseWriter, r *http.Request) {

}

// POST /pullRequest/reassign
func (*PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {

}
