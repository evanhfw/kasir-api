package domain

import "time"

type Transaction struct {
	ID          int                  `json:"id" example:"1"`
	TotalAmount int                  `json:"total_amount" example:"35000"`
	CreatedAt   time.Time            `json:"created_at" example:"2024-01-15T10:30:00Z"`
	Details     []TransactionDetails `json:"details"`
}

type TransactionDetails struct {
	ID            int    `json:"id" example:"1"`
	TransactionID int    `json:"transaction_id" example:"1"`
	ProductID     int    `json:"product_id" example:"5"`
	ProductName   string `json:"product_name" example:"Indomie Goreng"`
	Quantity      int    `json:"quantity" example:"3"`
	Subtotal      int    `json:"subtotal" example:"10500"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id" example:"5"`
	Quantity  int `json:"quantity" example:"3"`
}