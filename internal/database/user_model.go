package database

import "time"

type UserDB struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"Password"`
	CreatedAt time.Time `db:created_at`
}
