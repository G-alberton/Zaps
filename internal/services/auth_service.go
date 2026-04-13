package services

import (
	"ZAPS/internal/auth"
	"ZAPS/internal/repository"
)

type AuthService struct {
	Repo repository.UserRespository
	JWT  *auth.JWTService
}
