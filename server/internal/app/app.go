package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/robertt3kuk/tasks-golang/server/config"
	v1 "github.com/robertt3kuk/tasks-golang/server/internal/controller/http/v1"
	"github.com/robertt3kuk/tasks-golang/server/internal/service"
	"github.com/robertt3kuk/tasks-golang/server/internal/service/repository"
	"github.com/robertt3kuk/tasks-golang/server/pkg/httpserver"
	"github.com/robertt3kuk/tasks-golang/server/pkg/logger"
	"github.com/robertt3kuk/tasks-golang/server/pkg/mongodb"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	db := mongodb.New(cfg.Mongo.URL, cfg.Mongo.Name)
	mong := repository.New(db)
	service := service.New(mong)
	handler := chi.NewRouter()
	v1.NewRouter(handler, l, service)
	l.Info("app - Run - Started on port")
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
