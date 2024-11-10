package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"golang-apigw-lambdas/pkg/repository"

	"github.com/aws/aws-lambda-go/events"
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

func (h *UserHandler) HandleRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	users := []repository.User{
		{Name: "John Doe"},
		{Name: "Jane Doe"},
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: "Failed to parse users to JSON",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonUsers),
	}, nil
}
