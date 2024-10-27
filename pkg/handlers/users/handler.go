package handlers

import (
	"encoding/json"
	"net/http"

	"golang-apigw-lambdas/pkg/repository"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := h.repo.List()
	jsonUsers, err := json.Marshal(users)
	if err != nil {
		// err handling
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonUsers)
}

func (h *UserHandler) Register(r *chi.Mux) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.List)
	})
}
