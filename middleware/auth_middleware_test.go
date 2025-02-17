package middleware_test

import (
	"ecommerce/middleware"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock TokenVerifier
type MockVerifier struct {
	ValidToken string
	Err        error
}

func (m MockVerifier) VerifyToken(tokenString string) (string, error) {
	if tokenString == m.ValidToken {
		return "testuser", nil
	}
	return "", m.Err
}

func TestAuthMiddleware(t *testing.T) {
	mockVerifier := MockVerifier{
		ValidToken: "valid-token",
		Err:        errors.New("invalid token"),
	}

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		handler := middleware.Auth(mockVerifier, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		handler := middleware.Auth(mockVerifier, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		handler := middleware.Auth(mockVerifier, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
		}
	})
}
