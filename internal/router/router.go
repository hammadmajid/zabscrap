package router

import (
	"net/http"
	"zabscrap/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(application *app.App) *chi.Mux {
	router := chi.NewRouter()

	// Health check endpoint
	router.Get("/health", application.Handler.Health)

	// API endpoints
	router.Post("/fetch", application.Handler.FetchAttendance)
	router.Get("/api/build-info", application.Handler.BuildInfo)

	// Serve static files from web/ directory
	fileServer := http.FileServer(http.Dir("web"))
	router.Handle("/web/*", http.StripPrefix("/web", fileServer))

	// Serve index.html for root path
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	return router
}
