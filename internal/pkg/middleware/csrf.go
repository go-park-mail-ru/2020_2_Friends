package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/session"
	log "github.com/friends/pkg/logger"
)

type CSRFChecker struct {
	authChecker   AuthChecker
	sessionClient session.SessionWorkerClient
}

func NewCSRFChecker(authChecker AuthChecker, sessionClient session.SessionWorkerClient) CSRFChecker {
	return CSRFChecker{
		authChecker:   authChecker,
		sessionClient: sessionClient,
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

		userID, ok := r.Context().Value(UserID(configs.UserID)).(string)
		if !ok {
			err = fmt.Errorf("couldn't get userID from context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token, err := c.sessionClient.GetTokenFromUser(context.Background(), &session.UserID{Id: userID})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		tokenFromHeader := r.Header.Get("X-CSRF-Token")

		if token.GetValue() != tokenFromHeader {
			err = fmt.Errorf("token doesn't fit")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})

	return c.authChecker.Check(handlerFunc)
}
