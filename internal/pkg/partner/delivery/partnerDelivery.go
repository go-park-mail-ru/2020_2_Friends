package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	"github.com/friends/internal/pkg/vendors"
	log "github.com/friends/pkg/logger"
)

type PartnerDelivery struct {
	userUsecase    user.Usecase
	sessionUsecase session.Usecase
	vendorUsecase  vendors.Usecase
}

func New(userUsecase user.Usecase, sessionUsecase session.Usecase, vendorUsecase vendors.Usecase) PartnerDelivery {
	return PartnerDelivery{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
		vendorUsecase:  vendorUsecase,
	}
}

func (p PartnerDelivery) Create(w http.ResponseWriter, r *http.Request) {
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

	err = p.userUsecase.CheckIfUserExists(*user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	user.Role = configs.AdminRole

	userID, err := p.userUsecase.Create(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionName, err := p.sessionUsecase.Create(userID)
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
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusCreated)
}

func (p PartnerDelivery) CreateVendor(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	vendor := models.Vendor{}
	err = json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.vendorUsecase.Create(vendor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
