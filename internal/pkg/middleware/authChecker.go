package middleware

import (
	"context"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/session"
	log "github.com/friends/pkg/logger"
)

type UserID string

type AuthChecker struct {
	sessionClient session.SessionWorkerClient
}

func NewAuthChecker(sessionClient session.SessionWorkerClient) AuthChecker {
	return AuthChecker{
		sessionClient: sessionClient,
	}
}

func (a AuthChecker) Check(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		userID, err := a.sessionClient.Check(context.Background(), &session.SessionName{Name: cookie.Value})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserID(configs.UserID), userID.GetId())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
