package middleware

import (
	"fmt"
	"net/http"

	"github.com/friends/configs"
	log "github.com/friends/pkg/logger"
)

type CSRFChecker struct {
	authChecker AuthChecker
}

func NewCSRFChecker(authChecker AuthChecker) CSRFChecker {
	return CSRFChecker{
		authChecker: authChecker,
	}
}

func (c CSRFChecker) Check(next http.HandlerFunc) http.Handler {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				log.ErrorLogWithCtx(r.Context(), err)
			}
		}()

		cookie, err := r.Cookie(configs.CookieCSRF)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		tokenFromCookie := cookie.Value
		tokenFromHeader := r.Header.Get("X-CSRF-Token")

		if tokenFromCookie != tokenFromHeader {
			err = fmt.Errorf("token doesn't fit")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})

	return c.authChecker.Check(handlerFunc)
}
