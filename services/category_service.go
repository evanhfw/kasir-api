package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	r *repositories.CategoryRepository
}

func NewCategoryService(r *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{r: r}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.r.GetAll()
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.r.Create(category)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.r.GetByID(id)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.r.Update(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.r.Delete(id)
}