package services

import (
	"task-3/dto"
	"task-3/models"
	"task-3/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(name string) ([]dto.ProductResponse, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) GetByID(id int) (dto.ProductResponse, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
