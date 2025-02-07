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

func TestGetByID(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock query response
	rows := sqlmock.NewRows([]string{"id", "name", "email", "username", "password"}).
		AddRow(1, "Digambar", "diga123@gmail.com", "diga123", "diga@123")

	mock.ExpectQuery("select id, name, email, username, password from users where id=?").
		WithArgs(1).WillReturnRows(rows)

	userRepo := repository.NewUserRepo(db)

	user, err := userRepo.GetByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, "Digambar", user.Name)
	assert.Equal(t, "diga123@gmail.com", user.Email)
	assert.Equal(t, "diga123", user.Username)
	assert.Equal(t, "diga@123", user.Password)
}
