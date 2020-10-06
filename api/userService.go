package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/friends/model"
	"github.com/friends/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db storage.DB
}

func (us UserService) login(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println("error")
		return
	}

	ok := us.db.CheckUser(nil, user.Login, user.Password)
	if !ok {
		fmt.Println("user %v not auth", user)
		w.Write([]byte(`{"authorized": false}`))
		return
	}
	fmt.Println("user %v auth", user)
	expiration := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "testcookie",
		Expires:  expiration,
		Secure:   true,
		SameSite: 4,
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte(`{"authorized": true}`))
}

func (us UserService) reginster(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println("error")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		return
	}

	err = us.db.CreateUser(nil, user.Login, string(hashedPassword))
	if err != nil {
		fmt.Println("user %v not reg", user)
		w.Write([]byte(`{"created": false}`))

		return
	}
	fmt.Println("user %v reg", user)
	w.Write([]byte(`{"created": true}`))
}

func (us UserService) testCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("no cookie", err)
		return
	}
	fmt.Println("users cookie: ", cookie)
}
