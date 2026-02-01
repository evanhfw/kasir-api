package service

import (
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

// CategoryService handles category business logic
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(category *domain.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) GetByID(id int) (*domain.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(category *domain.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
