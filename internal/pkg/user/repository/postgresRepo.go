package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
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

func (ur UserRepository) Create(user models.User) (userID string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't hash password: %w", err)
	}

	err = ur.db.QueryRow(
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		user.Login, hashedPassword,
	).Scan(&userID)

	if err != nil {
		return "", fmt.Errorf("couldn't insert user in Postgres: %w", err)
	}

	return userID, nil
}

func (ur UserRepository) CheckIfUserExists(user models.User) error {
	row := ur.db.QueryRow(
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

func (ur UserRepository) CheckLoginAndPassword(user models.User) (userID string, err error) {
	row := ur.db.QueryRow(
		"SELECT id, password FROM users WHERE login=$1",
		user.Login,
	)

	dbUser := models.User{}
	switch err := row.Scan(&dbUser.ID, &dbUser.Password); err {
	case sql.ErrNoRows:
		return "", fmt.Errorf("user doesn't exist")

	case nil:
		isEqual := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if isEqual != nil {
			return "", fmt.Errorf("wrong password: %w", isEqual)
		}
		return dbUser.ID, nil

	default:
		return "", fmt.Errorf("db error: %w", err)
	}
}

func (u UserRepository) Delete(userID string) error {
	id, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("couldn't convert to string: %w", err)
	}
	_, err = u.db.Exec(
		"DELETE FROM users WHERE id=$1",
		id,
	)

	if err != nil {
		return fmt.Errorf("couldn't delete user from Postgres: %w", err)
	}

	return nil
}
