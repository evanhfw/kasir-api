package repository

import (
	"database/sql"
	"time"

	"kasir-api/internal/domain"
)

type reportRepository struct {
	db *sql.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{db: db}
}

// GetSalesSummary retrieves sales summary for the given date range
func (r *reportRepository) GetSalesSummary(startDate, endDate time.Time) (*domain.SalesSummary, error) {
	summary := &domain.SalesSummary{}

	// Get total revenue and transaction count
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*) 
		FROM transactions 
		WHERE created_at >= $1 AND created_at < $2
	`, startDate, endDate).Scan(&summary.TotalRevenue, &summary.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// Get best selling product
	var productName string
	var quantitySold int
	err = r.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&productName, &quantitySold)

	if err == sql.ErrNoRows {
		// No transactions in this period, best_selling_product will be null
		summary.BestSellingProduct = nil
	} else if err != nil {
		return nil, err
	} else {
		summary.BestSellingProduct = &domain.BestSellingProduct{
			Name:         productName,
			QuantitySold: quantitySold,
		}
	}

	return summary, nil
}
