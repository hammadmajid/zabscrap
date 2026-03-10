package app

import (
	"log"
	"zabscrap/internal/api"
)

type App struct {
	Logger  *log.Logger
	Handler api.Handler
}

func NewApp(logger *log.Logger, handler api.Handler) *App {
	return &App{
		Logger:  logger,
		Handler: handler,
	}
}
