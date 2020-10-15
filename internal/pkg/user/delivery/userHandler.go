package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/user"
)

type UserHandler struct {
	userUsecase user.Usecase
}

func NewUserHandler(usecase user.Usecase) UserHandler {
	return UserHandler{
		userUsecase: usecase,
	}
}

func (uh UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uh.userUsecase.Create(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (uh UserHandler) Verify(user models.User) (userID string, err error) {
	return uh.userUsecase.Verify(user)
}
