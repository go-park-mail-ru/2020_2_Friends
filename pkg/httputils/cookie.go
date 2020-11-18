package httputils

import (
	"net/http"
	"time"

	"github.com/friends/configs"
)

func SetCookie(w http.ResponseWriter, cookieValue string) {
	expiration := time.Now().Add(configs.ExpireTime)
	cookie := http.Cookie{
		Name:     configs.SessionID,
		Value:    cookieValue,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func DeleteCookie(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}

func SetCookieWithSameSiteNone(w http.ResponseWriter, cookieValue string) {
	expiration := time.Now().Add(configs.ExpireTime)
	cookie := http.Cookie{
		Name:     configs.SessionID,
		Value:    cookieValue,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
}

func SetAdminsCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(configs.ExpireTime)
	cookie := http.Cookie{
		Name:    configs.AdminsCookieName,
		Value:   "true",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
}

func SetAdminsCookieWithSameSiteNone(w http.ResponseWriter) {
	expiration := time.Now().Add(configs.ExpireTime)
	cookie := http.Cookie{
		Name:     configs.AdminsCookieName,
		Value:    "true",
		Expires:  expiration,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
}
