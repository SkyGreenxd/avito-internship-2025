package server

import (
	"avito-internship/pkg/logger"
	"context"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort         = "8080"
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultKeepAlive    = 60 * time.Second
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	KeepAlive    time.Duration
}

// Server обёртка над http.Server для запуска и остановки HTTP-сервиса.
type Server struct {
	httpServer *http.Server
}

// LoadHttpServerConfig загружает конфигурацию HTTP-сервера из переменных окружения
func LoadHttpServerConfig(logger logger.Logger) Config {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		logger.Warnf("the environment variable HTTP_PORT is not set. Using default value %s.", defaultPort)
		port = defaultPort
	}

	return Config{
		Port:         port,
		ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", defaultReadTimeout, logger),
		WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout, logger),
		KeepAlive:    getEnvAsDuration("KEEP_ALIVE", defaultKeepAlive, logger),
	}
}

// NewServer создаёт новый HTTP-сервер с заданным обработчиком и конфигурацией.
func NewServer(handler http.Handler, httpServer Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + httpServer.Port,
			Handler:      handler,
			ReadTimeout:  httpServer.ReadTimeout,
			WriteTimeout: httpServer.WriteTimeout,
			IdleTimeout:  httpServer.KeepAlive,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// getEnvAsDuration вспомогательная функция для парсинга duration
func getEnvAsDuration(key string, fallback time.Duration, logger logger.Logger) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		logger.Warnf("the environment variable %s is not set. Using default value %s.", key, fallback)
		return fallback
	}

	duration, err := time.ParseDuration(valStr)
	if err != nil {
		logger.Warnf("invalid duration value for %s: '%s'. Using default %v\n", key, valStr, fallback)
		return fallback
	}

	return duration
}
