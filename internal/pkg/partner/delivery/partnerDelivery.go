package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	"github.com/friends/internal/pkg/vendors"
	log "github.com/friends/pkg/logger"
	"github.com/gorilla/mux"
)

type PartnerDelivery struct {
	userUsecase    user.Usecase
	sessionUsecase session.Usecase
	vendorUsecase  vendors.Usecase
	profileUsecase profile.Usecase
}

func New(userUsecase user.Usecase, profileUsecase profile.Usecase, sessionUsecase session.Usecase, vendorUsecase vendors.Usecase) PartnerDelivery {
	return PartnerDelivery{
		userUsecase:    userUsecase,
		profileUsecase: profileUsecase,
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

	err = p.profileUsecase.Create(userID)
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

func (p PartnerDelivery) UpdateVendor(w http.ResponseWriter, r *http.Request) {
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

	vendorID := mux.Vars(r)["id"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vendor := models.Vendor{}
	err = json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.vendorUsecase.Update(vendor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p PartnerDelivery) AddProductToVendor(w http.ResponseWriter, r *http.Request) {
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

	vendorID := mux.Vars(r)["id"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product := models.Product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.VendorID, err = strconv.Atoi(vendorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = p.vendorUsecase.AddProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p PartnerDelivery) UpdateVendorPicture(w http.ResponseWriter, r *http.Request) {
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

	vendorID := mux.Vars(r)["id"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(configs.ImgMaxSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	file, _, err := r.FormFile(configs.ImgFormFileKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgName, err := p.vendorUsecase.UpdateVendorPicture(vendorID, file)
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

func (p PartnerDelivery) UpdateProductOnVendor(w http.ResponseWriter, r *http.Request) {
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

	vendorID := mux.Vars(r)["id"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product := models.Product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.VendorID, err = strconv.Atoi(vendorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = p.vendorUsecase.UpdateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p PartnerDelivery) DeleteProductFromVendor(w http.ResponseWriter, r *http.Request) {
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

	vendorID := mux.Vars(r)["id"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product := models.Product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.vendorUsecase.DeleteProduct(product.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p PartnerDelivery) UpdateProductPicture(w http.ResponseWriter, r *http.Request) {
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

	productID := mux.Vars(r)["id"]

	vendorID, err := p.vendorUsecase.GetVendorIDFromProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(configs.ImgMaxSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	file, _, err := r.FormFile(configs.ImgFormFileKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgName, err := p.vendorUsecase.UpdateProductPicture(productID, file)
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
