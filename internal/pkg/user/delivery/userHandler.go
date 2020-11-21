package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	"github.com/friends/pkg/httputils"
	log "github.com/friends/pkg/logger"
)

type UserHandler struct {
	userUsecase    user.Usecase
	sessionUsecase session.Usecase
	profileUsecase profile.Usecase
}

func NewUserHandler(usecase user.Usecase, sessionUsecase session.Usecase, profileUsecase profile.Usecase) UserHandler {
	return UserHandler{
		userUsecase:    usecase,
		sessionUsecase: sessionUsecase,
		profileUsecase: profileUsecase,
	}
}

func (u UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	user := &models.User{}
	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.Sanitize()

	err = u.userUsecase.CheckIfUserExists(*user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	user.Role = configs.UserRole

	userID, err := u.userUsecase.Create(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = u.profileUsecase.Create(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionValue, err := u.sessionUsecase.Create(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.SetCookie(w, sessionValue)
	w.WriteHeader(http.StatusCreated)
}

func (u UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	userID, err := u.sessionUsecase.Check(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.userUsecase.Delete(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = u.profileUsecase.Delete(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = u.sessionUsecase.Delete(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.DeleteCookie(w, cookie)
}
