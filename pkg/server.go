package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HttpHandler interface {
	Register(r *chi.Mux)
}

type LambdaHandler interface {
	HandleRequest(ctx context.Context, event json.RawMessage) error
}

type Server struct {
	router   *chi.Mux
	config   *Config
	handlers []HttpHandler
}

type Config struct {
	Environment string
	Port        string
	DbUser      string
	DbPassword  string
	DbPort      string
	DbName      string
	DbHost      string
}

func LoadConfig() *Config {
	return &Config{
		Environment: GetEnvOrDefault("APP_ENV", "development"),
		Port:        GetEnvOrDefault("APP_PORT", "8080"),
		DbHost:      GetEnvOrDefault("DB_HOST", "localhost.localstack.cloud"),
		DbUser:      GetEnvOrDefault("DB_USER", "localhost"),
		DbPassword:  GetEnvOrDefault("DB_PASSWORD", "localhost"),
		DbPort:      GetEnvOrDefault("DB_PORT", "4510"),
		DbName:      GetEnvOrDefault("DB_NAME", "foundflix"),
	}
}

func (s *Server) LoadDBPool() (*pgxpool.Pool, error) {
	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		s.config.DbUser,
		s.config.DbPassword,
		s.config.DbHost,
		s.config.DbPort,
		s.config.DbName,
	)

	// configure connection pool
	config, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %v", err)
	}

	// customize pool settings
	config.MaxConns = 25                      // Maximum number of connections in the pool
	config.MinConns = 5                       // Minimum number of connections to maintain
	config.MaxConnLifetime = 1 * time.Hour    // Maximum lifetime of a connection
	config.MaxConnIdleTime = 30 * time.Minute // Maximum idle time before closin
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// ping the database
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return pool, nil
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

func (s *Server) RegisterHandlers(handlers ...HttpHandler) {
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
