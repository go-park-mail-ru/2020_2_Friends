package storage

import (
	"database/sql"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type MapDB struct {
	db map[string]string
	*sync.RWMutex
}

func NewMapDB() MapDB {
	return MapDB{
		make(map[string]string),
		&sync.RWMutex{},
	}
}

func (m MapDB) CreateUser(db *sql.DB, login string, password string) error {
	if _, ok := m.db[login]; ok {
		return fmt.Errorf("user exists")
	}

	m.Lock()
	m.db[login] = password
	m.Unlock()

	return nil
}

func (m MapDB) CheckUser(db *sql.DB, login string, password string) bool {
	m.RLock()
	pass, ok := m.db[login]
	m.RUnlock()

	isEqual := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	return ok && (isEqual == nil)
}
