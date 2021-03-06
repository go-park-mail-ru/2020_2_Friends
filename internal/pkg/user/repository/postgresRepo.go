package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
	ownErr "github.com/friends/pkg/error"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return UserRepository{
		db: db,
	}
}

func (u UserRepository) Create(user models.User) (userID string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't hash password: %w", err)
	}

	err = u.db.QueryRow(
		"INSERT INTO users (login, password, role) VALUES ($1, $2, $3) RETURNING id",
		user.Login, hashedPassword, user.Role,
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("couldn't insert user in Postgres: %w", err)
	}

	return userID, nil
}

func (u UserRepository) CheckIfUserExists(user models.User) error {
	row := u.db.QueryRow(
		"SELECT login FROM users WHERE login=$1",
		user.Login,
	)

	var login string
	switch err := row.Scan(&login); err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return fmt.Errorf("user exists")
	default:
		return fmt.Errorf("error with db: %w", err)
	}
}

func (u UserRepository) CheckLoginAndPassword(user models.User) (userID string, err error) {
	row := u.db.QueryRow(
		"SELECT id, password FROM users WHERE login=$1",
		user.Login,
	)

	dbUser := models.User{}
	switch err := row.Scan(&dbUser.ID, &dbUser.Password); err {
	case sql.ErrNoRows:
		return "", ownErr.NewClientError(err)

	case nil:
		bcryptErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if bcryptErr != nil {
			return "", ownErr.NewClientError(fmt.Errorf("wrong password: %w", bcryptErr))
		}
		return dbUser.ID, nil

	default:
		return "", ownErr.NewServerError(fmt.Errorf("db error: %w", err))
	}
}

func (u UserRepository) Delete(userID string) error {
	_, err := u.db.Exec(
		"DELETE FROM users WHERE id=$1",
		userID,
	)

	if err != nil {
		return fmt.Errorf("couldn't delete user from Postgres: %w", err)
	}

	return nil
}

func (u UserRepository) CheckUsersRole(userID string) (int, error) {
	row := u.db.QueryRow(
		"SELECT role from users WHERE id = $1",
		userID,
	)

	var role int

	err := row.Scan(&role)
	if err != nil {
		return 0, fmt.Errorf("couldn't get role from user: %w", err)
	}

	return role, nil
}
