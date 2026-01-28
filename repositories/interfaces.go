package repositories

import "kasir-api/models"

// ProductRepository defines the interface for product data access
type ProductRepositoryInterface interface {
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
}

// CategoryRepository defines the interface for category data access
type CategoryRepositoryInterface interface {
	GetAll() ([]models.Category, error)
	Create(category *models.Category) error
	GetByID(id int) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id int) error
}
