package repository_test

import (
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

func TestGetByUsername(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "username", "password"}).
		AddRow(1, "Digambar", "diga123@gmail.com", "diga123", "diga@123")

	mock.ExpectQuery("select id, name, email, username, password from users where username=?").
		WithArgs("diga123").WillReturnRows(rows)

	userRepo := repository.NewUserRepo(db)

	user, err := userRepo.GetByUsername("diga123")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, "Digambar", user.Name)
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "username", "password"}).
		AddRow(1, "Digambar", "diga123@gmail.com", "diga123", "diga@123").
		AddRow(2, "Yash", "yash123@gmail.com", "yash123", "yash@123")

	mock.ExpectQuery("select id, name, email, username, password from users").
		WillReturnRows(rows)

	userRepo := repository.NewUserRepo(db)

	users, err := userRepo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Digambar", users[0].Name)
	assert.Equal(t, "Yash", users[1].Name)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`(?i)UPDATE users SET name=\?, email=\?, username=\?, password=\? WHERE id=\?`).
		WithArgs("Updated Name", "updated@example.com", "updateduser", "newpassword", 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // ✅ Correct result for UPDATE query

	userRepo := repository.NewUserRepo(db)

	user := &models.User{
		Id:       1,
		Name:     "Updated Name",
		Email:    "updated@example.com",
		Username: "updateduser",
		Password: "newpassword",
	}

	err = userRepo.Update(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet()) // ✅ Ensure expectations are met
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("delete from users where id=?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := repository.NewUserRepo(db)

	err = userRepo.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
