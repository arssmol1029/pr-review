package handlers

import (
	"net/http"
	"pr-review/internal/service"
)

type PRHandler struct {
	service *service.PRService
}

func NewPRHandler(s *service.PRService) *PRHandler {
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
