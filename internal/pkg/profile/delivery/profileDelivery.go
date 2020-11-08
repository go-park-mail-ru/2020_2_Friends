package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
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

	userID, ok := r.Context().Value(middleware.UserID(configs.UserID)).(string)
	if !ok {
		err = fmt.Errorf("couldn't get userID from context")
		w.WriteHeader(http.StatusInternalServerError)
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

	userID, ok := r.Context().Value(middleware.UserID(configs.UserID)).(string)
	if !ok {
		err = fmt.Errorf("couldn't get userID from context")
		w.WriteHeader(http.StatusInternalServerError)
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

func (p ProfileDelivery) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	userID, ok := r.Context().Value(middleware.UserID(configs.UserID)).(string)
	if !ok {
		err = fmt.Errorf("couldn't get userID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(configs.ImgMaxSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	file, _, err := r.FormFile(configs.AvatarFormFileKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgName, err := p.profUsecase.UpdateAvatar(userID, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := models.ImgResponse{Avatar: imgName}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
