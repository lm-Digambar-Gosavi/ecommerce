package services

import (
	"errors"
	"testing"

	"ecommerce/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock // provides mocking functionality
}

func (m *MockUserRepo) Create(user *models.User) error {
	args := m.Called(user) // Calls the Create method
	return args.Error(0)   // returns the first argument
}

func (m *MockUserRepo) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	user := args.Get(0)
	if user != nil { // user object is found, it returns the user and the error
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1) // no user is found, it returns nil and the error.
}

func (m *MockUserRepo) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user := args.Get(0)
	if user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepo) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepo) // Creates a mock repository
	userService := NewUserService(mockRepo)

	user := &models.User{
		Id:       1,
		Username: "abhay123",
		Password: "abhay@123",
	}
	t.Run("success	", func(t *testing.T) {
		mockRepo.On("GetByUsername", "abhay123").Return(user, nil)

		token, err := userService.Login("abhay123", "abhay@123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token) // token should be generated.
		// Verify that all expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail (incorrect password)", func(t *testing.T) {
		mockRepo.On("GetByUsername", "abhay123").Return(user, nil)
		token, err := userService.Login("abhay123", "wrong_password")
		assert.Error(t, err) // should return an error
		assert.Empty(t, token)
		assert.Equal(t, "invalid username or password", err.Error())
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})

	t.Run("fail (Not exist user)", func(t *testing.T) {
		mockRepo.On("GetByUsername", "non_existent").Return(nil, errors.New("not found"))
		token, err := userService.Login("non_existent", "password")
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	user := &models.User{
		Id:       1,
		Name:     "Abhay",
		Email:    "abhay123@gmail.com",
		Username: "abhay123",
		Password: "abhay@123",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", 1).Return(nil, errors.New("not found"))
		mockRepo.On("Create", user).Return(nil)

		err := userService.CreateUser(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
	t.Run("Alredy exist", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", 1).Return(user, nil)

		err := userService.CreateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "id already registered", err.Error())
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
	t.Run("Missing Fields", func(t *testing.T) {
		tests := []struct {
			name string
			user *models.User
		}{
			{
				"Missing Name", &models.User{Id: 2, Email: "abhay123@gmail.com", Password: "abhay@123"},
			},
			{
				"Missing Email", &models.User{Id: 2, Name: "Abhay", Password: "abhay@123"},
			},
			{
				"Missing Password", &models.User{Id: 2, Name: "Abhay", Email: "abhay123@gmail.com"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := userService.CreateUser(tc.user)
				assert.Error(t, err)
				assert.Equal(t, "all fields are required", err.Error())
			})
		}
	})
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)
	mockuser := &models.User{
		Id:       1,
		Name:     "Abhay",
		Email:    "abhay123@gmail.com",
		Username: "abhay123",
		Password: "abhay@123",
	}

	t.Run("User found", func(t *testing.T) {
		mockRepo.On("GetByID", 1).Return(mockuser, nil)
		user, err := userService.GetUserByID(1)
		assert.NoError(t, err)
		assert.Equal(t, mockuser, user)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo.On("GetByID", 99).Return(nil, errors.New("not found"))
		user, err := userService.GetUserByID(99)
		assert.Error(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
}

func TestGetAllUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)

	mockUsers := []models.User{
		{
			Id:       1,
			Name:     "Abhay",
			Email:    "abhay123@gmail.com",
			Username: "abhay123",
			Password: "abhay@123",
		},
		{
			Id:       2,
			Name:     "Yash",
			Email:    "yash123@gmail.com",
			Username: "yash123",
			Password: "yash@123"},
	}

	t.Run("User found", func(t *testing.T) {
		mockRepo.On("GetAll").Return(mockUsers, nil)
		users, err := userService.GetAllUser()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(users))
		assert.Equal(t, mockUsers, users)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})

	t.Run("Not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetAll").Return([]models.User{}, errors.New("database error"))
		users, err := userService.GetAllUser()
		assert.Error(t, err)
		assert.Empty(t, users)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)

	user := &models.User{
		Id:       1,
		Name:     "Abhay",
		Email:    "abhay123@gmail.com",
		Password: "abhay@123",
	}

	t.Run("User found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", 1).Return(user, nil)
		mockRepo.On("Update", user).Return(nil)

		err := userService.UpdateUser(user)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", 99).Return(nil, nil) // Simulating user not found

		err := userService.UpdateUser(&models.User{
			Id:       99,
			Name:     "Abhay",
			Email:    "abhay123@gmail.com",
			Password: "abhay@123",
		})
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())

		mockRepo.AssertExpectations(t)
	})
	t.Run("Failure -Database Error", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("GetByID", 1).Return(nil, errors.New("database error"))

		err := userService.UpdateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error()) // Expected output matches function behavior

		mockRepo.AssertExpectations(t)
	})
	t.Run("Missing Fields", func(t *testing.T) {
		tests := []struct {
			name string
			user *models.User
		}{
			{
				"Missing Name", &models.User{Id: 2, Email: "abhay123@gmail.com", Password: "abhay@123"},
			},
			{
				"Missing Email", &models.User{Id: 2, Name: "Abhay", Password: "abhay@123"},
			},
			{
				"Missing Password", &models.User{Id: 2, Name: "Abhay", Email: "abhay123@gmail.com"},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				err := userService.UpdateUser(tc.user)
				assert.Error(t, err)
				assert.Equal(t, "all fields are required", err.Error())
			})
		}
	})
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo)

	t.Run("User found", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil)
		err := userService.DeleteUser(1)
		assert.NoError(t, err)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockRepo.On("Delete", 99).Return(errors.New("not found"))
		err := userService.DeleteUser(99)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t) // Verify that all expectations were met
	})
}
