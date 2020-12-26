package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	ownErr "github.com/friends/pkg/error"
	image "github.com/friends/pkg/image"
	log "github.com/friends/pkg/logger"
)

type ProfileDelivery struct {
	profUsecase profile.Usecase
}

func NewProfileDelivery(profUsecase profile.Usecase) ProfileDelivery {
	return ProfileDelivery{
		profUsecase: profUsecase,
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	profile.Sanitize()

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
	profile.Sanitize()
	log.DataLog(r, profile)

	profile.UserID = userID

	err = p.profUsecase.Update(*profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	file, header, err := r.FormFile(configs.AvatarFormFileKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	imageType, err := image.CheckMimeType(mimeType)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	imgName, err := p.profUsecase.UpdateAvatar(userID, file, imageType)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusUnsupportedMediaType)
		return
	}

	resp := models.ImgResponse{Avatar: imgName}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (p ProfileDelivery) UpdateAddresses(w http.ResponseWriter, r *http.Request) {
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

	profile := models.Profile{}
	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	profile.Sanitize()
	log.DataLog(r, profile)

	err = p.profUsecase.UpdateAddresses(userID, profile.Addresses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
