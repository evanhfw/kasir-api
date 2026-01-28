package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	r repositories.ProductRepositoryInterface
}

func NewProductService(r repositories.ProductRepositoryInterface) *ProductService {
	return &ProductService{r: r}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.r.GetAll()
}

func (s *ProductService) Create(product *models.Product) error {
	return s.r.Create(product)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.r.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.r.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.r.Delete(id)
}
