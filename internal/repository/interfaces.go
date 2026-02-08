package repository

import (
	"kasir-api/internal/domain"
	"time"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetAll(name string) ([]domain.Product, error)
	Create(product *domain.Product) error
	GetByID(id int) (*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id int) error
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	Create(category *domain.Category) error
	GetByID(id int) (*domain.Category, error)
	Update(category *domain.Category) error
	Delete(id int) error
}

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	GetSalesSummary(startDate, endDate time.Time) (*domain.SalesSummary, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	CreateTransaction(items []domain.CheckoutItem) (*domain.Transaction, error)
}
