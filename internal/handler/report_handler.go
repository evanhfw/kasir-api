package handler

import (
	"net/http"
	"time"

	"kasir-api/internal/service"
)

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	service *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetTodaySalesSummary godoc
// @Summary      Get today's sales summary
// @Description  Get sales summary for today
// @Tags         reports
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.SalesSummary
// @Failure      500  {object}  handler.APIResponse  "Failed to fetch sales summary"
// @Router       /report/hari-ini [get]
func (h *ReportHandler) GetTodaySalesSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Use zero time to trigger today's date in service
	summary, err := h.service.GetSalesSummary(time.Time{}, time.Time{})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch sales summary")
		return
	}

	WriteJSON(w, http.StatusOK, summary)
}

// GetSalesSummary godoc
// @Summary      Get sales summary by date range
// @Description  Get sales summary for a specific date range
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        start_date  query     string  true  "Start date (YYYY-MM-DD)"  example(2026-01-01)
// @Param        end_date    query     string  true  "End date (YYYY-MM-DD)"    example(2026-02-01)
// @Success      200         {object}  domain.SalesSummary
// @Failure      400         {object}  handler.APIResponse  "Invalid date format"
// @Failure      500         {object}  handler.APIResponse  "Failed to fetch sales summary"
// @Router       /report [get]
func (h *ReportHandler) GetSalesSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Validate required parameters
	if startDateStr == "" || endDateStr == "" {
		WriteError(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}
	// Add one day to include the end date fully
	endDate = endDate.AddDate(0, 0, 1)

	summary, err := h.service.GetSalesSummary(startDate, endDate)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch sales summary")
		return
	}

	WriteJSON(w, http.StatusOK, summary)
}
