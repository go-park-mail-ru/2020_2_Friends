package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	userDelivery "github.com/friends/internal/pkg/user/delivery"
)

type SessionDelivery struct {
	sessionUsecase session.Usecase
	userDelivery   user.Delivery
}

func NewSessionDelivery(usecase session.Usecase, userDelivery userDelivery.UserHandler) SessionDelivery {
	return SessionDelivery{
		sessionUsecase: usecase,
		userDelivery:   userDelivery,
	}
}

func (sd SessionDelivery) Create(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := sd.userDelivery.Verify(*user)
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

func (sd SessionDelivery) Check(sessionName string) (userID string, err error) {
	userID, err = sd.sessionUsecase.Check(sessionName)

	return userID, err
}

func (sd SessionDelivery) Delete(sessionName string) error {
	return sd.sessionUsecase.Delete(sessionName)
}
