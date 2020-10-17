package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	log "github.com/friends/pkg/logger"
)

type SessionDelivery struct {
	sessionUsecase session.Usecase
	userUsecase    user.Usecase
}

func NewSessionDelivery(usecase session.Usecase, userUsecase user.Usecase) SessionDelivery {
	return SessionDelivery{
		sessionUsecase: usecase,
		userUsecase:    userUsecase,
	}
}

func (sd SessionDelivery) Create(w http.ResponseWriter, r *http.Request) {
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

	userID, err := sd.userUsecase.Verify(*user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionName, err := sd.sessionUsecase.Create(userID)
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
}

func (sd SessionDelivery) Delete(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = sd.sessionUsecase.Delete(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
}
