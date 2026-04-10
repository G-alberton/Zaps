package repository

import (
	"ZAPS/internal/database"
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

func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	dbUser := &database.UserDB{}

	query := `
		SELECT id, name, email, password, created_at
        FROM users
        WHERE email = $1
	`

	err := r.DB.QueryRow(query, email).Scan(
		&dbUser.ID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.Password,
		&dbUser.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return mapper.ToDomain(dbUser), nil
}

func (r *PostgresUserRepository) FindByID(id int64) (*models.User, error) {
	dbUser := &database.UserDB{}

	query := `
        SELECT id, name, email, password, created_at
        FROM users
        WHERE id = $1
    `

	err := r.DB.QueryRow(query, id).Scan(
		&dbUser.ID,
		&dbUser.Name,
		&dbUser.Email,
		&dbUser.Password,
		&dbUser.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return mapper.ToDomain(dbUser), nil
}
