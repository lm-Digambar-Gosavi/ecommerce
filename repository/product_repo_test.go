package repository

import (
	"database/sql"
	"ecommerce/models"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	product := &models.Product{
		ID:    1,
		Name:  "TubeLight",
		Price: 999,
	}
	repo := NewProductRepo(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("insert into products").
			WithArgs(product.Name, product.Price).
			WillReturnResult(sqlmock.NewResult(1, 1))
		// 1 : The inserted row ID.
		// 1 : One row affected (successful insert).

		err = repo.Create(product)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail", func(t *testing.T) {
		mock.ExpectExec("insert into products").
			WithArgs(product.Name, product.Price).
			WillReturnError(fmt.Errorf("failed to insert product"))

		err = repo.Create(product)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestGetByIdProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepo(db)
	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery("select id, name, price from products where id=?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow(1, "TubeLight", 999))

		product, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 1, product.ID)
		assert.Equal(t, "TubeLight", product.Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("NotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("select id, name, price from products where id=?").
			WithArgs(90).
			WillReturnError(sql.ErrNoRows)

		repo := NewProductRepo(db)
		product, err := repo.GetByID(90)

		assert.Error(t, err)
		assert.Nil(t, product)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("fail", func(t *testing.T) {
		mock.ExpectQuery("select id, name, price from products where id=\\?").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		product, err := repo.GetByID(1)

		assert.Error(t, err) // Expect an error
		assert.Nil(t, product)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetAllProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewProductRepo(db)
	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("select id, name, price from products").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow(1, "TubeLight", 999).
				AddRow(2, "Laptop", 49999))

		products, err := repo.GetAll()

		assert.NoError(t, err)
		assert.Len(t, products, 2) // ensures that exactly 2 products were returned

		assert.Equal(t, 1, products[0].ID)
		assert.Equal(t, "TubeLight", products[0].Name)

		assert.Equal(t, 2, products[1].ID)
		assert.Equal(t, "Laptop", products[1].Name)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectQuery("select id, name, price from products").
			WillReturnError(fmt.Errorf("database error"))

		products, err := repo.GetAll()

		assert.Error(t, err)
		assert.Nil(t, products)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Scan Error", func(t *testing.T) {
		mock.ExpectQuery("select id, name, price from products").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}). // missing "password" column
										AddRow(1, "TV"))

		products, err := repo.GetAll()

		assert.Error(t, err)
		assert.Nil(t, products)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	product := &models.Product{
		ID:    1,
		Name:  "TubeLight",
		Price: 999,
	}
	repo := NewProductRepo(db)
	mock.ExpectExec("update products set name = \\?, price = \\? where id = \\?").
		WithArgs(product.Name, product.Price, product.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
		// Row ID = 1,
		// 1 row affected
	err = repo.Update(product)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	repo := NewProductRepo(db)

	t.Run("Found", func(t *testing.T) {
		mock.ExpectExec("delete from products where id=?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		// Row ID = 1,
		// 1 row affected
		err = repo.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectExec("delete from products where id=?").
			WithArgs(90).
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

		err = repo.Delete(90)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("Fail", func(t *testing.T) {
		mock.ExpectExec("delete from products where id=\\?").
			WithArgs(1).
			WillReturnError(fmt.Errorf("failed to delete product"))

		err := repo.Delete(1)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
