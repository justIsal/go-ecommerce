package service

import (
	"errors"
	"go-ecommerce/internal/domain"
	"go-ecommerce/internal/repository"
	"go-ecommerce/pkg/utils"
	"time"
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Register(input domain.User) error
	Login(email, password string) (*LoginResponse, error)
	Logout(refreshToken string) error
	Refresh(refreshToken string) (string, error) 
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(input domain.User) error {
	hashedPwd, err := utils.HashPassword(input.Password)
	if err != nil {
		return err
	}
	input.Password = hashedPwd
	
	if input.Role == "" {
		input.Role = domain.RoleUser
	}

	return s.repo.CreateUser(&input)
}

func (s *authService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	
	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refTokenModel := domain.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: utils.GetRefreshTokenExpiry(),
	}
	
	if err := s.repo.StoreRefreshToken(&refTokenModel); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *authService) Logout(refreshToken string) error {
	return s.repo.DeleteRefreshToken(refreshToken)
}

func (s *authService) Refresh(refreshTokenStr string) (string, error) {
	storedToken, err := s.repo.FindRefreshToken(refreshTokenStr)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		s.repo.DeleteRefreshToken(refreshTokenStr) 
		return "", errors.New("refresh token expired")
	}

	if storedToken.Revoked {
		return "", errors.New("token revoked")
	}
	
	user, err := s.repo.FindByID(storedToken.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Role)
	
	return newAccessToken, err
}