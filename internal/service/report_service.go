package service

import (
	"time"

	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

type ReportService struct {
	repo repository.ReportRepository
}

// NewReportService creates a new report service
func NewReportService(repo repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetSalesSummary retrieves sales summary for the given date range
// If startDate and endDate are zero, it returns today's summary
func (s *ReportService) GetSalesSummary(startDate, endDate time.Time) (*domain.SalesSummary, error) {
	// If no dates provided, use today
	if startDate.IsZero() || endDate.IsZero() {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 1)
	}

	return s.repo.GetSalesSummary(startDate, endDate)
}
