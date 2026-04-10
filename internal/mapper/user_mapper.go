package mapper

import (
	"ZAPS/internal/database"
	"ZAPS/internal/dto"
	"ZAPS/internal/models"
)

func ToDomain(u *database.UserDB) *models.User {
	return &models.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

func ToResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToDB(u *models.User) *database.UserDB {
	return &database.UserDB{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
