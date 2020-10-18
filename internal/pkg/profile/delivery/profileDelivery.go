package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/friends/internal/pkg/models"

	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/session"
	log "github.com/friends/pkg/logger"
)

type ProfileDelivery struct {
	profUsecase profile.Usecase
	sessUsecase session.Usecase
}

func NewProfileDelivery(pUsecase profile.Usecase, sUsecase session.Usecase) ProfileDelivery {
	return ProfileDelivery{
		profUsecase: pUsecase,
		sessUsecase: sUsecase,
	}
}

func (p ProfileDelivery) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := p.sessUsecase.Check(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile, err := p.profUsecase.Get(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p ProfileDelivery) Update(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := p.sessUsecase.Check(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile := &models.Profile{}
	err = json.NewDecoder(r.Body).Decode(profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profile.UserID = userID

	err = p.profUsecase.Update(*profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
