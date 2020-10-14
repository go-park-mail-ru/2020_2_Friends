package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
)

type UserHandler struct {
	UserUsecase    user.Usecase
	SessionHandler session.Delivery
}

func (uh UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uh.UserUsecase.Create(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
