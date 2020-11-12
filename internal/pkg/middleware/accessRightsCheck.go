package middleware

import (
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/user"
)

type AccessRightsChecker struct {
	userUsecase user.Usecase
}

func NewAccessRightsChecker(userUsecase user.Usecase) AccessRightsChecker {
	return AccessRightsChecker{
		userUsecase: userUsecase,
	}
}

func (a AccessRightsChecker) AccessRightsCheck(next http.HandlerFunc, neccessaryRole int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserID(configs.UserID)).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		role, err := a.userUsecase.CheckUsersRole(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if role != neccessaryRole {
			w.WriteHeader(http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}
