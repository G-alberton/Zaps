package repository

import "ZAPS/internal/models"

type UserRespository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
}
