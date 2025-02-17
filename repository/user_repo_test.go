package repository

import (
	"database/sql"
	"ecommerce/models"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	// db is a mocked database connection that behaves like a real *sql.DB obj
	// mock is the controller that helps us set expectations
	assert.NoError(t, err)
	defer db.Close()

	user := &models.User{ // sample user obj
		Id:       1,
		Name:     "Abhay",
		Email:    "abhay123@gmail.com",
		Username: "abhay123",
		Password: "abhay@123",
	}

	repo := NewUserRepo(db) // inject mock db

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("insert into users"). //expect an insert statement
							WithArgs(user.Name, user.Email, user.Username, user.Password). //expected arguments
							WillReturnResult(sqlmock.NewResult(1, 1))                      // returning a mock result.
		// 1 → The inserted row ID.
		// 1 → One row affected (successful insert).

		err = repo.Create(user)
		assert.NoError(t, err) // checks if the method returns an error.
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Failer", func(t *testing.T) {
		mock.ExpectExec("insert into users").
			WithArgs(user.Name, user.Email, user.Username, user.Password).
			WillReturnError(fmt.Errorf("failed to insert user"))

		err = repo.Create(user)
		assert.Error(t, err) // error due to failed query
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByIdUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users where id=?").
			WithArgs(1). // query should be called with id=1
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "username", "password"}).AddRow(1, "Abhay", "abhay123@gmail.com", "abhay123", "abhay@123"))

		user, err := repo.GetByID(1)

		assert.NoError(t, err) // should not return an error
		assert.NotNil(t, user) // returned user should not be nil
		assert.Equal(t, "abhay123", user.Username)
		assert.Equal(t, "abhay@123", user.Password)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Scan error", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users where id=?").
			WithArgs(1).                                                               // query should be called with id=1
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Abhay")) // Missing required columns

		user, err := repo.GetByID(1)
		assert.Error(t, err) // error due to scan failure
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users where id=?").
			WithArgs(90).
			WillReturnError(sql.ErrNoRows) // "database/sql"

		user, err := repo.GetByID(90)

		assert.Error(t, err) // function must return an error
		assert.Nil(t, user)  // user should be nil.

		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestGetByUsernameUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users where username=?").
			WithArgs("abhay123").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "username", "password"}).AddRow(1, "Abhay", "abhay123@gmail.com", "abhay123", "abhay@123"))

		user, err := repo.GetByUsername("abhay123")

		assert.NoError(t, err)
		assert.NotNil(t, user)

		assert.Equal(t, 1, user.Id)
		assert.Equal(t, "abhay@123", user.Password)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users where username=?").
			WithArgs("abc@123").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByUsername("abc@123")

		assert.Error(t, err)
		assert.Nil(t, user)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetAllUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "username", "password"}).
				AddRow(1, "Abhay", "abhay123@gmail.com", "abhay123", "abhay@123").
				AddRow(2, "Alesh", "alesh123@gmail.com", "alesh123", "alesh@123"))

		users, err := repo.GetAll()

		assert.NoError(t, err)
		assert.Len(t, users, 2) // ensures that exactly 2 users were returned

		assert.Equal(t, 1, users[0].Id)
		assert.Equal(t, "abhay123", users[0].Username)
		assert.Equal(t, "abhay@123", users[0].Password)

		assert.Equal(t, 2, users[1].Id)
		assert.Equal(t, "alesh123", users[1].Username)
		assert.Equal(t, "alesh@123", users[1].Password)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users").
			WillReturnError(fmt.Errorf("database error"))

		users, err := repo.GetAll()

		assert.Error(t, err)
		assert.Nil(t, users)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Scan Error", func(t *testing.T) {
		mock.ExpectQuery("select id, name, email, username, password from users").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "username"}). // missing "password" column
													AddRow(1, "Abhay", "abhay123@gmail.com", "abhay123"))

		users, err := repo.GetAll()

		assert.Error(t, err)
		assert.Nil(t, users)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	user := &models.User{
		Id:       1,
		Name:     "Abhay",
		Email:    "abhay123@gmail.com",
		Username: "abhay123",
		Password: "abhay@123",
	}

	mock.ExpectExec(regexp.QuoteMeta("update users set name=?, email=?, username=?, password=? where id=?")).
		WithArgs(user.Name, user.Email, user.Username, user.Password, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// Row ID = 1,
	// 1 row affected
	repo := NewUserRepo(db)
	err = repo.Update(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	t.Run("Found", func(t *testing.T) {
		mock.ExpectExec("delete from users where id=?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		// Row ID = 1,
		// 1 row affected

		err = repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectExec("delete from users where id=?").
			WithArgs(90).
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

		repo := NewUserRepo(db)
		err = repo.Delete(90)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail", func(t *testing.T) {
		mock.ExpectExec("delete from users where id=?").
			WithArgs(1).
			WillReturnError(fmt.Errorf("failed to delete user"))

		err = repo.Delete(1)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
