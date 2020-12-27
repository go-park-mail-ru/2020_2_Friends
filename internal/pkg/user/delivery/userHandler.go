package delivery

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	ownErr "github.com/friends/pkg/error"
	"github.com/friends/pkg/httputils"
	log "github.com/friends/pkg/logger"
)

type UserHandler struct {
	userUsecase    user.Usecase
	sessionClient  session.SessionWorkerClient
	profileUsecase profile.Usecase
}

func NewUserHandler(
	usecase user.Usecase, sessionClient session.SessionWorkerClient, profileUsecase profile.Usecase,
) UserHandler {
	return UserHandler{
		userUsecase:    usecase,
		sessionClient:  sessionClient,
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

	sessionName, err := u.sessionClient.Create(context.Background(), &session.UserID{Id: userID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.SetCookie(w, sessionName.GetName())

	token, err := u.sessionClient.SetCSRFToken(context.Background(), &session.UserID{Id: userID})
	httputils.SetCSRFCookie(w, token.GetValue())

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

	userID, err := u.sessionClient.Check(context.Background(), &session.SessionName{Name: cookie.Value})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.userUsecase.Delete(userID.GetId())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = u.profileUsecase.Delete(userID.GetId())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = u.sessionClient.Delete(context.Background(), &session.SessionName{Name: cookie.Value})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.DeleteCookie(w, cookie)
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	userID, err := u.userUsecase.Verify(*user)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	sessionName, err := u.sessionClient.Create(context.Background(), &session.UserID{Id: userID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.SetCookie(w, sessionName.GetName())

	token, err := u.sessionClient.SetCSRFToken(context.Background(), &session.UserID{Id: userID})
	httputils.SetCSRFCookie(w, token.GetValue())
}

func (u UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
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

	_, err = u.sessionClient.Delete(context.Background(), &session.SessionName{Name: cookie.Value})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	httputils.DeleteCookie(w, cookie)
}

func (u UserHandler) IsAuthorized(w http.ResponseWriter, r *http.Request) {
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

	_, err = u.sessionClient.Check(context.Background(), &session.SessionName{Name: cookie.Value})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
