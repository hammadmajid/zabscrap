package router

import (
	"zabscrap/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(application *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", application.Handler.Health)
	router.Get("/", application.Handler.ShowForm)
	router.Post("/fetch", application.Handler.FetchAttendance)

	return router
}
