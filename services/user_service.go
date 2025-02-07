package services

import (
	"ecommerce/models"
	"ecommerce/repository"
	"ecommerce/utils"
	"errors"
	"fmt"
)

type UserService interface {
	Login(username, password string) (string, error)
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetAllUser() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if user.Password != password {
		return "", errors.New("invalid username or password")
	}

	token, err := utils.CreateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) CreateUser(user *models.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("all fields are required")
	}

	existingUser, _ := s.userRepo.GetByID(user.Id)
	if existingUser != nil {
		return errors.New("id already registered")
	}

	fmt.Println("User registered successfully")
	return s.userRepo.Create(user)
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetAllUser() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) UpdateUser(user *models.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("all fields are required")
	}

	existingUser, err := s.userRepo.GetByID(user.Id)
	if err != nil || existingUser == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}
