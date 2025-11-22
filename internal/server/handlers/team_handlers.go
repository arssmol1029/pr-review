package handlers

import (
	"context"
	"net/http"
	"pr-review/internal/models"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

type TeamHandler struct {
	service TeamService
}

func NewTeamHandler(s TeamService) *TeamHandler {
	return &TeamHandler{service: s}
}

// POST /team/add
func (*TeamHandler) Add(w http.ResponseWriter, r *http.Request) {

}

// GET /team/get
func (*TeamHandler) Get(w http.ResponseWriter, r *http.Request) {

}
