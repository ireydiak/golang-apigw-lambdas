package main

import (
	handlers "golang-apigw-lambdas/pkg/handlers/users"
	"golang-apigw-lambdas/pkg/repository"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	repo := repository.NewUserRepo()
	handler := handlers.NewUserHandler(repo)
	lambda.Start(handler.HandleRequest)
}
