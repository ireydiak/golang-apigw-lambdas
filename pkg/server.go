package server

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler interface {
	Register(r *chi.Mux)
}

type Server struct {
	router   *chi.Mux
	config   *Config
	handlers []Handler
}

type Config struct {
	Environment string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		Environment: GetEnvOrDefault("APP_ENV", "development"),
		Port:        GetEnvOrDefault("APP_PORT", "8080"),
	}
}

func GetEnvOrDefault(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func NewServer(cfg *Config) *Server {
	router := chi.NewMux()
	router.Use(middleware.Logger)
	return &Server{
		router:   router,
		config:   cfg,
		handlers: nil,
	}
}

func (s *Server) RegisterHandlers(handlers ...Handler) {
	for _, h := range handlers {
		h.Register(s.router)
	}
	s.handlers = handlers
}

func (s *Server) Start() error {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		adapter := chiadapter.New(s.router)
		lambda.Start(adapter.ProxyWithContext)
		return nil
	}
	//return http.ListenAndServe(":"+s.config.Port, s.router)
	return http.ListenAndServe(":3000", s.router)
}
