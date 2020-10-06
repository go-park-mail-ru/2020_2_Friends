package storage

import (
	"database/sql"

	usr "github.com/friends/model/user"
)

type DB interface {
	CreateUser(db *sql.DB, login string, password string) error
	CheckUser(db *sql.DB, login string, password string) bool
	GetUser(db *sql.DB, login string) usr.User
	UpdateUser(db *sql.DB, user usr.User) error
	DeleteUser(db *sql.DB, login string)
}

type SessionDB interface {
	SaveCookie(db *sql.DB, cookie string)
	CheckCookie(db *sql.DB, cookie string) bool
	DeleteCookie(db *sql.DB, cookie string)
}
