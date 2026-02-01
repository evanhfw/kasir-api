package service

import (
	"errors"

	"kasir-api/internal/apperrors"
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetAll() ([]domain.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) Create(product *domain.Product) error {
	// Validate category exists
	_, err := s.categoryRepo.GetByID(product.CategoryID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrCategoryNotFound
		}
		return err
	}
	return s.productRepo.Create(product)
}

func (s *ProductService) GetByID(id int) (*domain.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) Update(product *domain.Product) error {
	// Validate category exists
	_, err := s.categoryRepo.GetByID(product.CategoryID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrCategoryNotFound
		}
		return err
	}
	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.productRepo.Delete(id)
}
