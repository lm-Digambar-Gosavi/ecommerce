package services

import (
	"errors"
	"testing"

	"ecommerce/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepo struct {
	mock.Mock // provides mocking functionality
}

func (m *MockProductRepo) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepo) GetByID(id int) (*models.Product, error) {
	args := m.Called(id)
	product := args.Get(0)
	if product != nil {
		return product.(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepo) GetAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepo) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepo)
	productService := NewProductService(mockRepo)
	validProduct := &models.Product{
		ID:    1,
		Name:  "Laptop",
		Price: 61000,
	}
	invalid := &models.Product{
		ID:    2,
		Name:  "Mouse",
		Price: 0,
	}
	t.Run("Valid Product", func(t *testing.T) {
		mockRepo.On("Create", validProduct).Return(nil)

		err := productService.CreateProduct(validProduct)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Invalid price", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		err := productService.CreateProduct(invalid)
		assert.Error(t, err)
		assert.Equal(t, "product price must be greter than zero", err.Error())
	})
}

func TestGetProductByID(t *testing.T) {
	mockRepo := new(MockProductRepo)
	productService := NewProductService(mockRepo)
	mockProduct := &models.Product{
		ID:    1,
		Name:  "Laptop",
		Price: 61000,
	}
	t.Run("Product Found", func(t *testing.T) {
		mockRepo.On("GetByID", 1).Return(mockProduct, nil)
		product, err := productService.GetProductByID(1)
		assert.NoError(t, err)
		assert.Equal(t, mockProduct, product)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Product not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", 90).Return(nil, errors.New("product not found"))
		product, err := productService.GetProductByID(90)
		assert.Error(t, err)
		assert.Nil(t, product)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllProduct(t *testing.T) {
	mockRepo := new(MockProductRepo)
	productService := NewProductService(mockRepo)
	mockProducts := []models.Product{
		{
			ID:    1,
			Name:  "Laptop",
			Price: 61000,
		},
		{
			ID:    2,
			Name:  "Mouse",
			Price: 0,
		},
	}
	t.Run("Product Found", func(t *testing.T) {
		mockRepo.On("GetAll").Return(mockProducts, nil)
		product, err := productService.GetAllProducts()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(product))
		assert.Equal(t, mockProducts, product)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
	t.Run("Not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll").Return([]models.Product{}, errors.New("database error"))
		product, err := productService.GetAllProducts()
		assert.Error(t, err)
		assert.Empty(t, product)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepo)
	productService := NewProductService(mockRepo)
	product := &models.Product{
		ID:    1,
		Name:  "Laptop",
		Price: 61000,
	}
	updatePro := &models.Product{
		ID:    1,
		Name:  "Gaming Laptop",
		Price: 61000,
	}
	invalidPro := &models.Product{
		ID:    1,
		Name:  "",
		Price: -61000,
	}
	t.Run("valid product", func(t *testing.T) {
		mockRepo.On("GetByID", 1).Return(product, nil)
		mockRepo.On("Update", updatePro).Return(nil)

		err := productService.UpdateProduct(updatePro)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Product not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", 90).Return(nil, errors.New("product not found"))

		err := productService.UpdateProduct(&models.Product{
			ID:    90,
			Name:  "New Product",
			Price: 1000,
		})
		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
	t.Run("Invalid product details", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		err := productService.UpdateProduct(invalidPro)
		assert.Error(t, err)
		assert.Equal(t, "all fields are required", err.Error())
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepo)
	productService := NewProductService(mockRepo)

	t.Run("Product deleted", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil)
		err := productService.DeleteProducts(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Product not found", func(t *testing.T) {
		mockRepo.On("Delete", 99).Return(errors.New("product not found"))
		err := productService.DeleteProducts(99)
		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
