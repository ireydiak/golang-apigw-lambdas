package main

import (
	server "golang-apigw-lambdas/pkg"
	handlers "golang-apigw-lambdas/pkg/handlers/users"
	"golang-apigw-lambdas/pkg/repository"
	"log"
)

func main() {
	cfg := server.LoadConfig()
	server := server.NewServer(cfg)

	repo := repository.NewUserRepo()
	handler := handlers.NewUserHandler(repo)
	server.RegisterHandlers(handler)

	log.Fatal(server.Start())
}
