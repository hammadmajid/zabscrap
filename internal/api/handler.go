package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

// BuildInfo returns git commit information from GitHub API
//
//goland:noinspection GoUnusedParameter
func (h *Handler) BuildInfo(w http.ResponseWriter, r *http.Request) {
	const githubAPI = "https://api.github.com/repos/hammadmajid/zabscrap/commits/main"

	var buildInfo struct {
		Hash      string `json:"hash"`
		Message   string `json:"message"`
		TimeAgo   string `json:"timeAgo"`
		Available bool   `json:"available"`
	}

	// Fetch latest commit from GitHub
	resp, err := http.Get(githubAPI)
	if err != nil {
		h.logger.Printf("Failed to fetch GitHub commit info: %v", err)
		buildInfo.Available = false
		h.sendJSON(w, http.StatusOK, JSONResponse{
			Success: true,
			Data:    buildInfo,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.logger.Printf("GitHub API returned status: %d", resp.StatusCode)
		buildInfo.Available = false
		h.sendJSON(w, http.StatusOK, JSONResponse{
			Success: true,
			Data:    buildInfo,
		})
		return
	}

	var githubResponse struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Date string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubResponse); err != nil {
		h.logger.Printf("Failed to decode GitHub response: %v", err)
		buildInfo.Available = false
		h.sendJSON(w, http.StatusOK, JSONResponse{
			Success: true,
			Data:    buildInfo,
		})
		return
	}

	// Extract short hash (first 7 characters)
	hash := githubResponse.SHA
	if len(hash) > 7 {
		hash = hash[:7]
	}

	// Parse commit time
	commitTime, err := time.Parse(time.RFC3339, githubResponse.Commit.Author.Date)
	if err != nil {
		h.logger.Printf("Failed to parse commit time: %v", err)
		buildInfo.Available = false
		h.sendJSON(w, http.StatusOK, JSONResponse{
			Success: true,
			Data:    buildInfo,
		})
		return
	}

	// Calculate time ago
	duration := time.Since(commitTime)
	var timeAgo string

	if duration.Hours() < 1 {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			timeAgo = "1 minute ago"
		} else {
			timeAgo = fmt.Sprintf("%d minutes ago", minutes)
		}
	} else if duration.Hours() < 24 {
		hours := int(duration.Hours())
		if hours == 1 {
			timeAgo = "1 hour ago"
		} else {
			timeAgo = fmt.Sprintf("%d hours ago", hours)
		}
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			timeAgo = "1 day ago"
		} else {
			timeAgo = fmt.Sprintf("%d days ago", days)
		}
	}

	buildInfo.Hash = hash
	buildInfo.Message = githubResponse.Commit.Message
	buildInfo.TimeAgo = timeAgo
	buildInfo.Available = true

	h.sendJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Data:    buildInfo,
	})
}
