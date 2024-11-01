package main

import (
	handlers "golang-apigw-lambdas/pkg/handlers/movies"
	"golang-apigw-lambdas/pkg/repository"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	//cfg := server.LoadConfig()
	//server := server.NewServer(cfg)
	//dbPool, err := server.LoadDBPool()
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}

	var dbPool *pgxpool.Pool
	repo := repository.NewMovieRepo(dbPool)
	handler := handlers.NewMoviesHandler(repo)
	lambda.Start(handler.HandleRequest)
	// server.RegisterHandlers(handler)
	//
	// log.Fatal(server.Start())
}
