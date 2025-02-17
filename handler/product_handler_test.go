package handler

import (
	"bytes"
	"context"
	"ecommerce/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ProductService

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) GetProductByID(id int) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductService) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProducts(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateProductHandler(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHander(mockService)

	product := models.Product{
		ID:    1,
		Name:  "Mouse",
		Price: 900,
	}
	body, _ := json.Marshal(product)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("Post", "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockService.On("CreateProduct", &product).Return(nil)
		handler.CreateProduct(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, res.Body.String(), "Product created successfully")
	})
	t.Run("Invalid request", func(t *testing.T) {
		product := `{
		"name": "Mouse",
		"price": "nine hundred"
		}`
		req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(product))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.CreateProduct(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request")
	})
	t.Run("Fail", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockService.On("CreateProduct", &product).Return(errors.New("database error"))

		handler.CreateProduct(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestGetProductByID(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHander(mockService)
	product := &models.Product{
		ID:    1,
		Name:  "Mouse",
		Price: 999,
	}

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetProductByID", 1).Return(product, nil)

		r := chi.NewRouter()                            // Create Chi router
		r.Get("/products/{id}", handler.GetProductByID) // Create request

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req) // Serve the request

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("invalid product id", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/products/{id}", handler.GetProductByID)

		req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid product ID")
	})
	t.Run("fail", func(t *testing.T) {
		mockService.On("GetProductByID", 99).Return(nil, errors.New("Product not found"))

		r := chi.NewRouter()
		r.Get("/products/{id}", handler.GetProductByID)

		req := httptest.NewRequest(http.MethodGet, "/products/99", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "Product not found")
	})
}

func TestGetAllProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHander(mockService)

	products := []models.Product{
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
	t.Run("Success", func(t *testing.T) {
		mockService.On("GetAllProducts").Return(products, nil)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("Get", "/products", nil)

		handler.GetAllProducts(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), `"Name":"Laptop"`)
		assert.Contains(t, res.Body.String(), `"Name":"Mouse"`)

	})
	t.Run("Fail", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAllProducts").Return([]models.Product{}, errors.New("Failed to retrieve products"))

		req := httptest.NewRequest("Get", "/products", nil)
		res := httptest.NewRecorder()

		handler.GetAllProducts(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "Failed to retrieve products")
	})
}

func TestUpdateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHander(mockService)

	product := models.Product{
		ID:    1,
		Name:  "Laptop",
		Price: 61000,
	}
	reqBody, _ := json.Marshal(product)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("Put", "/update", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockService.On("UpdateProduct", &product).Return(nil)

		handler.UpdateProduct(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Product updated successfully")
	})
	t.Run("Invalid request", func(t *testing.T) {
		product := `{
			"ID":    1,
			"Name":  "Laptop",
			"Price": ,
		}`

		res := httptest.NewRecorder()
		req := httptest.NewRequest("Put", "/update", bytes.NewBuffer([]byte(product)))

		handler.UpdateProduct(res, req)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Contains(t, res.Body.String(), "Invalid request ")
	})
	t.Run("Fail", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		res := httptest.NewRecorder()
		req := httptest.NewRequest("Put", "/update", bytes.NewBuffer(reqBody))

		mockService.On("UpdateProduct", &product).Return(errors.New("Failed"))

		handler.UpdateProduct(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
	})
}

func TestDeleteProducts(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHander(mockService)

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeleteProducts", 1).Return(nil)

		req := httptest.NewRequest("DELETE", "/products/1", nil)
		res := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext() //creates a new chi router context
		chiCtx.URLParams.Add("id", "1")

		// Attach Route Context to Request
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		handler.DeleteProducts(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Product deleted successfully")
	})
	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/products/invalid", nil)
		res := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext() //creates a new chi router context
		chiCtx.URLParams.Add("id", "invalid")

		// Attach Route Context to Request
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		handler.DeleteProducts(res, req)

		assert.Equal(t, res.Code, http.StatusBadRequest)
		assert.Contains(t, res.Body.String(), "Invalid Product ID")
	})
	t.Run("Fail", func(t *testing.T) {
		mockService.On("DeleteProducts", 90).Return(errors.New("Failed to delete product"))

		req := httptest.NewRequest("DELETE", "/products/90", nil)
		res := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext() //creates a new chi router context
		chiCtx.URLParams.Add("id", "90")

		// Attach Route Context to Request
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		handler.DeleteProducts(res, req)

		assert.Equal(t, res.Code, http.StatusInternalServerError)
		assert.Contains(t, res.Body.String(), "Failed to delete product")
	})
}
