package services

import (
	"ecommerce/models"
	"ecommerce/repository"
	"fmt"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id int) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProducts(id int) error
}

type productService struct {
	productRepo repository.ProductRepo
}

func NewProductService(productRepo repository.ProductRepo) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(product *models.Product) error {
	if product.Price <= 0 {
		return fmt.Errorf("product price must be greter than zero")
	}
	return s.productRepo.Create(product)
}

func (s *productService) GetProductByID(id int) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *productService) UpdateProduct(product *models.Product) error {
	if product.Name == "" || product.Price <= 0 {
		return fmt.Errorf("all fields are required")
	}

	existingProduct, err := s.productRepo.GetByID(product.ID)
	if err != nil || existingProduct == nil {
		return fmt.Errorf("product not found")
	}

	return s.productRepo.Update(product)
}

func (s *productService) DeleteProducts(id int) error {
	return s.productRepo.Delete(id)
}
