package handlers

import (
	"context"
	"encoding/json"
	"golang-apigw-lambdas/pkg/repository"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-chi/chi/v5"
)

type MoviesHandler struct {
	repo *repository.MovieRepo
}

func NewMoviesHandler(repo *repository.MovieRepo) *MoviesHandler {
	return &MoviesHandler{
		repo: repo,
	}
}

func (h *MoviesHandler) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	movies, err := h.repo.List()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	jsonMovies, err := json.Marshal(movies)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("failed to parse movies to a json encoded string"))
		return
	}

	w.WriteHeader(200)
	w.Write(jsonMovies)
	return
}

func (h *MoviesHandler) HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	movies, err := h.repo.List()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	jsonMovies, err := json.Marshal(movies)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonMovies),
	}, nil
}

func (h *MoviesHandler) Register(r *chi.Mux) {
	r.Route("/movies", func(r chi.Router) {
		r.Get("/", h.List)
	})
}
