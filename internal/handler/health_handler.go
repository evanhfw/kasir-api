package handler

import (
	"database/sql"
	"net/http"
)

// healthStatus represents the health check response
type healthStatus struct {
	Status   string `json:"status" example:"healthy"`
	Database string `json:"database" example:"connected"`
}

// HealthHandler handles health check requests
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health godoc
// @Summary      Health check
// @Description  Check if the server and database are running
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthStatus
// @Failure      503  {object}  handler.APIResponse
// @Router       /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Check database connection
	if err := h.db.Ping(); err != nil {
		WriteError(w, http.StatusServiceUnavailable, "Database connection failed")
		return
	}

	WriteJSON(w, http.StatusOK, healthStatus{
		Status:   "healthy",
		Database: "connected",
	})
}
