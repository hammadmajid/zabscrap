package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"zabscrap/internal/models"
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

// renderTemplate renders a template with the given data
func (h *Handler) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	if err := h.templates.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Data":     data,
		"Template": templateName,
	}); err != nil {
		h.logger.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

// ShowForm displays the login form
func (h *Handler) ShowForm(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "layout", nil); err != nil {
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

	h.renderResultsTemplate(w, data)
}

// renderResultsTemplate renders the results template with attendance data
func (h *Handler) renderResultsTemplate(w http.ResponseWriter, data []models.CourseAttendance) {
	if err := h.templates.ExecuteTemplate(w, "layout", data); err != nil {
		h.logger.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
