package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
)

type UserHandler struct {
	userUsecase    user.Usecase
	sessionUsecase session.Usecase
}

func NewUserHandler(usecase user.Usecase, sessionUsecase session.Usecase) UserHandler {
	return UserHandler{
		userUsecase:    usecase,
		sessionUsecase: sessionUsecase,
	}
}

func (uh UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isExists := uh.userUsecase.CheckIfUserExists(*user)
	if isExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := uh.userUsecase.Create(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionName, err := uh.sessionUsecase.Create(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(configs.ExpireTime)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionName,
		Expires:  expiration,
		Secure:   true,
		SameSite: 4,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusCreated)
}
