package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"zabscrap/internal/scraper"
)

type Handler struct {
	logger    *log.Logger
	scraper   *scraper.Scraper
	templates *template.Template
}

func NewHandler(logger *log.Logger) *Handler {
	h := &Handler{
		logger:  logger,
		scraper: scraper.New(),
	}
	h.loadTemplates()
	return h
}

// loadTemplates loads all template files from internal/templates
func (h *Handler) loadTemplates() {
	var err error
	h.templates, err = template.ParseGlob(filepath.Join("internal/templates", "*.tmpl"))
	if err != nil {
		h.logger.Fatalf("Failed to load templates: %v", err)
	}
}

// Health returns the health status
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Status is available")
	if err != nil {
		h.logger.Println(err)
		return
	}
}

// ShowForm displays the login form wrapped in layout
func (h *Handler) ShowForm(w http.ResponseWriter, r *http.Request) {
	// Execute form template
	if err := h.templates.ExecuteTemplate(w, "form", nil); err != nil {
		h.logger.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// FetchAttendance handles the form submission and displays results
func (h *Handler) FetchAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	data, err := h.scraper.ScrapeAttendance(username, password)
	if err != nil {
		h.logger.Printf("Scraping error: %v", err)
		http.Error(w, "Failed to fetch attendance data", http.StatusInternalServerError)
		return
	}

	// Execute results template with data
	if err := h.templates.ExecuteTemplate(w, "results", data); err != nil {
		h.logger.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
