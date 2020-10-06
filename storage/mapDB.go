package storage

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/friends/model/vendor"
	"golang.org/x/crypto/bcrypt"
)

type UserMapDB struct {
	db map[string]string
	*sync.RWMutex
}

func NewUserMapDB() UserMapDB {
	return UserMapDB{
		make(map[string]string),
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
