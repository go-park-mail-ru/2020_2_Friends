package repository

import (
	"database/sql"
	"fmt"

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
		return "", err
	}

	err = ur.db.QueryRow(
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		user.Login, hashedPassword,
	).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (ur UserRepository) CheckIfUserExists(user models.User) bool {
	row := ur.db.QueryRow(
		"SELECT login FROM users WHERE login=$1",
		user.Login,
	)

	var login string
	switch err := row.Scan(&login); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		fmt.Println(err)
		return true
	}
}

func (ur UserRepository) CheckLoginAndPassword(user models.User) (userID string, err error) {
	row := ur.db.QueryRow(
		"SELECT * FROM users WHERE login=$1",
		user.Login,
	)

	dbUser := models.User{}
	switch err := row.Scan(&dbUser.ID, &dbUser.Login, &dbUser.Password); err {
	case sql.ErrNoRows:
		return "", fmt.Errorf("there is no such user")

	case nil:
		isEqual := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if isEqual != nil {
			return "", fmt.Errorf("wrong password")
		}
		return dbUser.ID, nil

	default:
		return "", fmt.Errorf("there is no such user")
	}
}
