package repository

import (
	"database/sql"
	"ecommerce/models"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ProductRepo interface {
	Create(product *models.Product) error
	GetByID(id int) (*models.Product, error)
	GetAll() ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product *models.Product) error {
	query := "insert into products (Name,Price) values (?,?)"
	_, err := r.db.Exec(query, product.Name, product.Price)
	if err != nil {
		return fmt.Errorf("failed to insert product: %v", err)
	}
	return nil
}

func (r *productRepo) GetByID(id int) (*models.Product, error) {
	query := "select id, name, price from products where id=?"
	row := r.db.QueryRow(query, id)

	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepo) GetAll() ([]models.Product, error) {
	query := "select id, name, price from products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepo) Update(product *models.Product) error {
	_, err := r.db.Exec("update products set name = ?, price = ? where id = ?",
		product.Name, product.Price, product.ID)
	return err
}

func (r *productRepo) Delete(id int) error {
	query := "delete from products where id=?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}
