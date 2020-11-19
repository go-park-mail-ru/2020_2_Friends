package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	ownErr "github.com/friends/pkg/error"
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
	user.Sanitize()

	userID, err := sd.userUsecase.Verify(*user)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	role, err := sd.userUsecase.CheckUsersRole(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionName, err := sd.sessionUsecase.Create(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(configs.ExpireTime)
	cookie := http.Cookie{
		Name:     configs.SessionID,
		Value:    sessionName,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)

	if role == configs.AdminRole {
		adminCookie := http.Cookie{
			Name:     configs.AdminsCookieName,
			Value:    "true",
			Expires:  expiration,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, &adminCookie)
	}
}

func (sd SessionDelivery) Delete(w http.ResponseWriter, r *http.Request) {
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

	err = sd.sessionUsecase.Delete(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	cookie.Path = "/"
	http.SetCookie(w, cookie)

	adminsCookie, err := r.Cookie(configs.AdminsCookieName)
	if err != nil {
		return
	}

	adminsCookie.Expires = time.Now().AddDate(0, 0, -1)
	adminsCookie.Path = "/"
	http.SetCookie(w, adminsCookie)
}
