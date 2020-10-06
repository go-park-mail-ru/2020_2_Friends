package storage

import (
	"database/sql"
	"fmt"
	"sync"

	usr "github.com/friends/model/user"
	"github.com/friends/model/vendor"
	"golang.org/x/crypto/bcrypt"
)

type UserMapDB struct {
	db     map[string]string
	userDB map[string]usr.User
	*sync.RWMutex
}

func NewUserMapDB() UserMapDB {
	return UserMapDB{
		make(map[string]string),
		make(map[string]usr.User),
		&sync.RWMutex{},
	}
}

func (m UserMapDB) CreateUser(db *sql.DB, login string, password string) error {
	if _, ok := m.db[login]; ok {
		return fmt.Errorf("user exists")
	}

	m.Lock()
	m.db[login] = password
	m.Unlock()

	return nil
}

func (m UserMapDB) CheckUser(db *sql.DB, login string, password string) bool {
	m.RLock()
	pass, ok := m.db[login]
	m.RUnlock()

	isEqual := bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
	return ok && (isEqual == nil)
}

func (m UserMapDB) GetUser(db *sql.DB, login string) usr.User {
	m.RLock()
	u, ok := m.userDB[login]
	m.Unlock()
	if !ok {
		return usr.User{
			Login: login,
		}
	}

	return u
}

func (m UserMapDB) UpdateUser(db *sql.DB, user usr.User) error {
	m.Lock()
	m.userDB[user.Login] = user
	m.Unlock()

	return nil
}

func (m UserMapDB) DeleteUser(db *sql.DB, login string) {
	m.Lock()
	delete(m.userDB, login)
	delete(m.db, login)
	m.Unlock()
}

type VendorMapDB struct{}

func (vm VendorMapDB) GetVendor(id string) (vendor.Vendor, error) {
	return vendor.Vendor{
		Name: "Veggie shop",
		Products: []vendor.Product{
			{
				PicturePath: "assets/vegan.png",
				Name:        "Toffee",
				Price:       "179 Р",
			},
			{
				PicturePath: "assets/vegan.png",
				Name:        "Toffee",
				Price:       "179 Р",
			},
		},
	}, nil
}

type SessionMapDB struct {
	db map[string]struct{}
	*sync.RWMutex
}

func NewSessionMapDB() SessionMapDB {
	return SessionMapDB{
		make(map[string]struct{}),
		&sync.RWMutex{},
	}
}

func (sm SessionMapDB) SaveCookie(db *sql.DB, cookie string) {
	sm.Lock()
	sm.db[cookie] = struct{}{}
	sm.Unlock()
}

func (sm SessionMapDB) CheckCookie(db *sql.DB, cookie string) bool {
	sm.RLock()
	_, ok := sm.db[cookie]
	sm.RUnlock()
	return ok
}

func (sm SessionMapDB) DeleteCookie(db *sql.DB, cookie string) {
	sm.Lock()
	delete(sm.db, cookie)
	sm.Unlock()
}
