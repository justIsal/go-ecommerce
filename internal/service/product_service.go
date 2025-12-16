package service

import (
	"errors"
	"go-ecommerce/internal/domain"
	"go-ecommerce/internal/repository"
)

type ProductService interface {
	GetAll() ([]domain.Product, error)
	GetByID(id uint) (*domain.Product, error)
	Create(input domain.Product) error
	Update(id uint, input domain.Product) error
	Delete(id uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAll() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) GetByID(id uint) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) Create(input domain.Product) error {
	if input.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return s.repo.Create(&input)
}

func (s *productService) Update(id uint, input domain.Product) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Name = input.Name
	existing.Description = input.Description
	existing.Price = input.Price
	existing.Stock = input.Stock
	existing.ImageURL = input.ImageURL

	return s.repo.Update(existing)
}

func (s *productService) Delete(id uint) error {
	return s.repo.Delete(id)
}