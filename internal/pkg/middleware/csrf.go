package middleware

import (
	"fmt"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/csrf"
	log "github.com/friends/pkg/logger"
)

type CSRFChecker struct {
	csrfUsecase csrf.Usecase
	authChecker AuthChecker
}

func NewCSRFChecker(authChecker AuthChecker, csrfUsecase csrf.Usecase) CSRFChecker {
	return CSRFChecker{
		authChecker: authChecker,
		csrfUsecase: csrfUsecase,
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

		cookie, err := r.Cookie(configs.SessionID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session := cookie.Value
		tokenFromHeader := r.Header.Get("X-CSRF-Token")

		isSame := c.csrfUsecase.Check(tokenFromHeader, session)
		if !isSame {
			err = fmt.Errorf("token doesn't fit")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})

	return c.authChecker.Check(handlerFunc)
}
