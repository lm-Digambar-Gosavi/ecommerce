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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(id int) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAllUser() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestLoginHandler(t *testing.T) {
	mockService := new(MockUserService) // Create mock service
	handler := NewUserHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		user := models.User{
			Username: "abhay123",
			Password: "abhay@123",
		}
		body, _ := json.Marshal(user)
		// httptest.NewRequest(method, url, body) => Creates a fake HTTP request
		req := httptest.NewRequest("Post", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder() // Captures the response from the handler

		mockService.On("Login", "abhay123", "abhay@123").Return("tokenString", nil) // sets up expectations
		handler.LoginHandler(res, req)                                              // call actual Handler

		assert.Equal(t, http.StatusOK, res.Code) // check status code 200

		var resp map[string]string              // store key-value pairs from the JSON response
		json.Unmarshal(res.Body.Bytes(), &resp) // decode the JSON response
		// .Bytes() extracts the raw response as a byte slice
		// json.Unmarshal converts the byte(JSON) into a Go data structure.
		assert.Equal(t, "tokenString", resp["token"])
	})
	t.Run("Invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.LoginHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request")
	})
	t.Run("Invalid Credentials", func(t *testing.T) {
		user := map[string]string{
			"Username": "abhay123",
			"Password": "abhay123",
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockService.On("Login", "abhay123", "abhay123").Return("", errors.New("invalid username or password"))

		handler.LoginHandler(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRegisterUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("success", func(t *testing.T) {
		user := models.User{
			Id:       1,
			Name:     "Abhay",
			Email:    "abhay123@gmail.com",
			Username: "abhay123",
			Password: "abhay@123",
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("Post", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockService.On("CreateUser", &user).Return(nil)

		handler.RegisterUser(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		var resp map[string]string
		json.Unmarshal(res.Body.Bytes(), &resp)
		assert.Equal(t, "User registered successfully", resp["message"])
		mockService.AssertExpectations(t)
	})
	t.Run("Invalid Request Body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.RegisterUser(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid request")
	})
	t.Run("Failure", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		user := models.User{
			Id:       1,
			Name:     "Abhay",
			Email:    "abhay123@gmail.com",
			Username: "abhay123",
			Password: "abhay@123",
		}
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		mockService.On("CreateUser", &user).Return(errors.New("failed to create user"))

		handler.RegisterUser(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Contains(t, res.Body.String(), "failed to create user")
		mockService.AssertExpectations(t)
	})

}
func TestGetUserByID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	user := &models.User{
		Id:       1,
		Name:     "Abhay",
		Email:    "abhay123@gmail.com",
		Username: "abhay123",
		Password: "abhay@123",
	}

	r := chi.NewRouter()                      // Create Chi router
	r.Get("/users/{id}", handler.GetUserByID) // Create request

	t.Run("Success", func(t *testing.T) {
		mockService.On("GetUserByID", 1).Return(user, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req) // Serve the request

		assert.Equal(t, http.StatusOK, res.Code)

		var resp models.User
		json.Unmarshal(res.Body.Bytes(), &resp)
		assert.Equal(t, user, &resp)
		mockService.AssertExpectations(t)
	})
	t.Run("Fail", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abhay", nil) // Invalid ID
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Contains(t, res.Body.String(), "Invalid user ID")
		mockService.AssertExpectations(t)

	})
	t.Run("Fail (Not Found)", func(t *testing.T) {
		mockService.On("GetUserByID", 99).Return((*models.User)(nil), errors.New("user not found"))

		req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, res.Body.String(), "User not found")

		mockService.AssertExpectations(t)
	})
}

func TestGetAllUsers(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	users := []models.User{
		{
			Id:       1,
			Name:     "Abhay",
			Email:    "abhay123@gmail.com",
			Username: "abhay123",
			Password: "abhay@123"},
		{
			Id:       2,
			Name:     "Yash",
			Email:    "yash123@gmail.com",
			Username: "yash123",
			Password: "yash@123",
		},
	}
	t.Run("Success", func(t *testing.T) {
		mockService.On("GetAllUser").Return(users, nil)

		req := httptest.NewRequest("GET", "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAllUsers(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []models.User
		json.Unmarshal(rec.Body.Bytes(), &resp)

		assert.Equal(t, users, resp)
		mockService.AssertExpectations(t)
	})
	t.Run("Empty user", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAllUser").Return([]models.User{}, nil)

		req := httptest.NewRequest("GET", "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAllUsers(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []models.User
		json.Unmarshal(rec.Body.Bytes(), &resp)

		assert.Empty(t, resp) // Verify response is empty
		mockService.AssertExpectations(t)
	})
	t.Run("Fail", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("GetAllUser").Return([]models.User{}, errors.New("database error")) // Return empty slice

		req := httptest.NewRequest("GET", "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAllUsers(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code) // Expecting 500
	})
}

func TestUpdateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	user := models.User{
		Id:       1,
		Name:     "Yash",
		Email:    "yash123@gmail.com",
		Username: "yash123",
		Password: "yash@123",
	}
	reqBody, _ := json.Marshal(user)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/user/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Inject the id parameter into the request context
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		mockService.On("GetUserByID", 1).Return(&user, nil)
		mockService.On("UpdateUser", &user).Return(nil)

		handler.UpdateUser(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "User updated successfully")
	})
	t.Run("Invalid User ID", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/user/abc", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc") // Non-numeric ID
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.UpdateUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid user ID")
	})
	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{
		"Id": 1, 
		"Name": "Yash", 
		"Email": "invalid-email",
		}`
		req := httptest.NewRequest("PUT", "/user/1", bytes.NewBuffer([]byte(invalidJSON)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.UpdateUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request")
	})

	t.Run("Failure - User Not Found", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		req := httptest.NewRequest("PUT", "/user/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		mockService.On("GetUserByID", 1).Return((*models.User)(nil), errors.New("user not found"))
		mockService.On("UpdateUser", &user).Return(errors.New("database error"))

		handler.UpdateUser(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "User not found")
	})
	t.Run("UpdateUser Failure", func(t *testing.T) {
		mockService.ExpectedCalls = nil // Clear previous expectations

		req := httptest.NewRequest("PUT", "/user/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		mockService.On("GetUserByID", 1).Return(&user, nil)
		mockService.On("UpdateUser", &user).Return(errors.New("database error"))

		handler.UpdateUser(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "database error")
	})
}

func TestDeleteUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeleteUser", 1).Return(nil)

		req := httptest.NewRequest("DELETE", "/users/1", nil)
		res := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext() //creates a new chi router context
		chiCtx.URLParams.Add("id", "1")

		// Attach Route Context to Request
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.DeleteUser(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Invalid User ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/users/abc", nil) // Non-numeric ID
		rec := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.DeleteUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid user ID")
	})
	t.Run("Delete Failure", func(t *testing.T) {
		mockService.On("DeleteUser", 2).Return(errors.New("user not found"))

		req := httptest.NewRequest("DELETE", "/users/2", nil)
		rec := httptest.NewRecorder()

		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		handler.DeleteUser(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to delete user")
	})
}
