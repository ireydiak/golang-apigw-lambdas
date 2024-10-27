package main

import (
	"fmt"
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

	fmt.Printf("before")
	log.Fatal(server.Start())
	fmt.Printf("here")
}
