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
	sessUsecase session.Usecase
}

func NewAuthChecker(sessUsecase session.Usecase) AuthChecker {
	return AuthChecker{
		sessUsecase: sessUsecase,
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID, err := a.sessUsecase.Check(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserID(configs.UserID), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
