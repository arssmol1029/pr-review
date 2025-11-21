package handlers

import (
	"net/http"
	"pr-review/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// POST /users/setIsActive
func (*UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {

}

// GET /users/getReview
func (*UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {

}
