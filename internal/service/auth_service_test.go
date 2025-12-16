package service_test

import (
	"go-ecommerce/internal/domain"
	"go-ecommerce/internal/service"
	"go-ecommerce/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) StoreRefreshToken(token *domain.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockUserRepo) FindRefreshToken(token string) (*domain.RefreshToken, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockUserRepo) DeleteRefreshToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	authService := service.NewAuthService(mockRepo)

	email := "test@example.com"
	password := "rahasia123"
	hashedPassword, _ := utils.HashPassword(password)

	dummyUser := &domain.User{
		ID:       1,
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
	}

	mockRepo.On("FindByEmail", email).Return(dummyUser, nil)
	
	mockRepo.On("StoreRefreshToken", mock.Anything).Return(nil)

	response, err := authService.Login(email, password)

	assert.NoError(t, err)             
	assert.NotNil(t, response)         
	assert.NotEmpty(t, response.AccessToken) 
	assert.NotEmpty(t, response.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	authService := service.NewAuthService(mockRepo)

	email := "test@example.com"
	hashedPassword, _ := utils.HashPassword("passwordAsli") 

	dummyUser := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}

	mockRepo.On("FindByEmail", email).Return(dummyUser, nil)

	response, err := authService.Login(email, "passwordSalah")

	assert.Error(t, err)                 
	assert.Nil(t, response)              
	assert.Equal(t, "invalid email or password", err.Error()) 
}