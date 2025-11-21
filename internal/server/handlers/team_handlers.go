package handlers

import (
	"net/http"
	"pr-review/internal/service"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(s *service.TeamService) *TeamHandler {
	return &TeamHandler{service: s}
}

// POST /team/add
func (*TeamHandler) Add(w http.ResponseWriter, r *http.Request) {

}

// GET /team/get
func (*TeamHandler) Get(w http.ResponseWriter, r *http.Request) {

}
