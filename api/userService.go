package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/friends/model"
	"github.com/friends/storage"
)

type UserService struct {
	db storage.DB
}

func (us UserService) login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	user := &model.User{}
	json.Unmarshal(body, user)

	ok := us.db.CheckUser(nil, user.Login, hashPassword(user.Password))
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	user := &model.User{}
	json.Unmarshal(body, user)

	err = us.db.CreateUser(nil, user.Login, hashPassword(user.Password))
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

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return string(hash[:])
}
