package repository

import (
	"errors"
	"go-ecommerce/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	
	StoreRefreshToken(token *domain.RefreshToken) error
	FindRefreshToken(tokenString string) (*domain.RefreshToken, error)
	DeleteRefreshToken(tokenString string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) StoreRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *userRepository) FindRefreshToken(tokenString string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := r.db.Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		return nil, errors.New("refresh token not found")
	}
	return &token, nil
}

func (r *userRepository) DeleteRefreshToken(tokenString string) error {
	return r.db.Where("token = ?", tokenString).Delete(&domain.RefreshToken{}).Error
}