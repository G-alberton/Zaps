package services

import (
	"ZAPS/internal/auth"
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
	"errors"
)

type AuthService struct {
	Repo repository.UserRespository
	JWT  *auth.JWTService
}

func (s *AuthService) Register(name, email, password string) (string, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hash,
	}

	err = s.Repo.Create(user)
	if err != nil {
		return "", err
	}

	token, err := s.JWT.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := auth.CheckPassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.JWT.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
