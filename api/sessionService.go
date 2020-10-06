package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	user "github.com/friends/model/user"
	"github.com/friends/storage"
)

type SessionService struct {
	userService UserService
	db          storage.SessionDB
}

func (ss SessionService) login(w http.ResponseWriter, r *http.Request) {
	user := &user.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println("error")
		return
	}

	ok := ss.userService.CheckUser(*user)
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
	ss.db.SaveCookie(nil, cookie.Value)
	w.Write([]byte(`{"authorized": true}`))
}

func (ss SessionService) logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `no sess`, 401)
		return
	}
	fmt.Println(session.Value)

	if ok := ss.db.CheckCookie(nil, session.Value); !ok {
		http.Error(w, `no sess`, 401)
		return
	}

	ss.db.DeleteCookie(nil, session.Value)
	fmt.Println(ss.db.CheckCookie(nil, session.Value), "это куки")

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}
