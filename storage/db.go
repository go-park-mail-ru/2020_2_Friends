package storage

import "database/sql"

type DB interface {
	CreateUser(db *sql.DB, login string, password string) error
	CheckUser(db *sql.DB, login string, password string) bool
}
