package repository_test

import (
	//"database/sql"
	"ecommerce/models"
	"ecommerce/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Define expected SQL query and arguments
	mock.ExpectExec("insert into users").
		WithArgs("Abhay Sonawane", "abhay123@gmail.com", "abha123", "abha@123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Initialize repository with mock DB
	userRepo := repository.NewUserRepo(db)

	// Test data
	user := &models.User{
		Name:     "Abhay Sonawane",
		Email:    "abhay123@gmail.com",
		Username: "abha123",
		Password: "abha@123",
	}

	// Call the method and check for errors
	err = userRepo.Create(user)
	assert.NoError(t, err)

	// Ensure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}
