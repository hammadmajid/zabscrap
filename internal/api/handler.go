package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zabscrap/internal/scraper"
)

type Handler struct {
	logger  *log.Logger
	scraper *scraper.Scraper
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger:  logger,
		scraper: scraper.New(),
	}
}

// JSONResponse represents the API response structure
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// sendJSON sends a JSON response
func (h *Handler) sendJSON(w http.ResponseWriter, statusCode int, response JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding JSON: %v", err)
	}
}

// sendError sends an error JSON response
func (h *Handler) sendError(w http.ResponseWriter, statusCode int, message string) {
	h.sendJSON(w, statusCode, JSONResponse{
		Success: false,
		Error:   message,
	})
}

// Health returns the health status
//
//goland:noinspection GoUnusedParameter
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Status is available")
	if err != nil {
		h.logger.Println(err)
		return
	}
}

// FetchAttendance handles the JSON API request for fetching attendance
func (h *Handler) FetchAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse JSON request body
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if reqBody.Username == "" || reqBody.Password == "" {
		h.sendError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	data, err := h.scraper.ScrapeAttendance(reqBody.Username, reqBody.Password)
	if err != nil {
		h.logger.Printf("Scraping error: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to fetch attendance data")
		return
	}

	h.sendJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    data,
	})
}
