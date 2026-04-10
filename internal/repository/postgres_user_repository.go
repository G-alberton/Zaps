package repository

import (
	"ZAPS/internal/mapper"
	"ZAPS/internal/models"
	"database/sql"
)

type PostgresUserRepository struct {
	Db *sql.DB
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	dbUser := mapper.ToDB(user)

	query := `
	INSERT INTO users (name, email, password)
    VALUES ($1, $2, $3)
    RETURNING id, created_at
`

	err := r.DB.QueryRow(
		query,
		dbUser.Name,
		dbUser.Email,
		dbUser.Password,
	).Scan(&dbUser.ID, &dbUser.CreatedAt)

	if err != nil {
		return err
	}

	user.ID = dbUser.ID

	return nil
}
