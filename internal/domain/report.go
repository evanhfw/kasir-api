package domain

// SalesSummary represents the sales summary report
// @Description Sales summary information
type SalesSummary struct {
	TotalRevenue       int                 `json:"total_revenue" example:"45000"`
	TotalTransactions  int                 `json:"total_transactions" example:"5"`
	BestSellingProduct *BestSellingProduct `json:"best_selling_product,omitempty"`
}

// BestSellingProduct represents the top selling product
// @Description Best selling product information
type BestSellingProduct struct {
	Name         string `json:"name" example:"Indomie Goreng"`
	QuantitySold int    `json:"quantity_sold" example:"12"`
}
