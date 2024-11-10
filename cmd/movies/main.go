package main

import (
	handlers "golang-apigw-lambdas/pkg/handlers/movies"
	"golang-apigw-lambdas/pkg/repository"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	repo := repository.NewMovieRepo()
	handler := handlers.NewMoviesHandler(repo)
	lambda.Start(handler.HandleRequest)
}
