package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Register(r chi.Router)
}

type Server struct {
	router   chi.Router
	config   *Config
	handlers []Handler
}

type Config struct {
	Environment string
	Port        string
}

func loadConfig() *Config {
	return &Config{
		Environment: "development",
		Port:        "8080",
	}
}

func NewServer(cfg *Config) *Server {
	router := chi.NewRouter()
	return &Server{
		router:   router,
		config:   cfg,
		handlers: nil,
	}
}

func (s *Server) RegisterHandlers(handlers ...Handler) {
	for _, h := range handlers {
		s.handlers = append(s.handlers, h)
	}
}

func (s *Server) Start() error {
	return http.ListenAndServe(":"+s.config.Port, s.router)
}

func main() {
	cfg := loadConfig()
	server := NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("%v", err)
	}
	// server.RegisterHandlers()

}
