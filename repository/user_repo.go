package repository

import (
	"database/sql"
	"ecommerce/models"
	"fmt"
)

type UserRepo interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	query := "insert into users (name, email, username, password) values (?,?,?,?)"
	_, err := r.db.Exec(query, user.Name, user.Email, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

func (r *userRepo) GetByID(id int) (*models.User, error) {
	query := "select id, name, email, username, password from users where id=?"
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	query := "select id, name, email, username, password from users where username=?"
	row := r.db.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *userRepo) GetAll() ([]models.User, error) {
	query := "select id, name, email, username, password from users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) Update(user *models.User) error {
	_, err := r.db.Exec("update users set name = ?, email = ?, username = ?, password = ? WHERE id = ?",
		user.Name, user.Email, user.Username, user.Password, user.Id)
	return err
}

func (r *userRepo) Delete(id int) error {
	query := "delete from users where id=?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user : %v", err)
	}
	return nil
}
