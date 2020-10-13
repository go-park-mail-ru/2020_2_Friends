package repository

import (
	"database/sql"

	"github.com/friends/internal/pkg/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur UserRepository) Create(user models.User) error {
	_, err := ur.db.Exec(
		"INSERT INTO users (login, password, name, email) VALUES ($1, $2, $3, $4)",
		user.Login, user.Password, user.Name, user.Email,
	)

	return err
}
