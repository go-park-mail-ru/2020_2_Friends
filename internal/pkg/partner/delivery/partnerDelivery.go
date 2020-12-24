package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/session"
	"github.com/friends/internal/pkg/user"
	"github.com/friends/internal/pkg/vendors"
	"github.com/friends/pkg/httputils"
	"github.com/friends/pkg/image"
	log "github.com/friends/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
)

type PartnerDelivery struct {
	userUsecase    user.Usecase
	sessionClient  session.SessionWorkerClient
	vendorUsecase  vendors.Usecase
	profileUsecase profile.Usecase
}

func New(
	userUsecase user.Usecase, profileUsecase profile.Usecase,
	sessionClient session.SessionWorkerClient, vendorUsecase vendors.Usecase,
) PartnerDelivery {
	return PartnerDelivery{
		userUsecase:    userUsecase,
		profileUsecase: profileUsecase,
		sessionClient:  sessionClient,
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

	session, err := p.sessionClient.Create(context.Background(), &session.UserID{Id: userID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputils.SetCookie(w, session.GetName())

	token := shortuuid.NewWithNamespace(session.GetName())
	httputils.SetCSRFCookie(w, token)

	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
	w.Header().Set("X-CSRF-Token", token)

	w.WriteHeader(http.StatusCreated)
}

func (p PartnerDelivery) CreateVendor(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	partnerID, ok := r.Context().Value(middleware.UserID(configs.UserID)).(string)
	if !ok {
		err = fmt.Errorf("couldn't get userID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vendor := models.Vendor{}
	err = json.NewDecoder(r.Body).Decode(&vendor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(vendor.Categories) > 5 {
		err = fmt.Errorf("too much categories for vendor")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vendorID, err := p.vendorUsecase.Create(partnerID, vendor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := models.AddResponse{
		ID: vendorID,
	}

	err = json.NewEncoder(w).Encode(resp)
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

	productID, err := p.vendorUsecase.AddProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := models.AddResponse{
		ID: productID,
	}

	err = json.NewEncoder(w).Encode(resp)
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

	file, header, err := r.FormFile(configs.ImgFormFileKey)
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

	imgName, err := p.vendorUsecase.UpdateVendorPicture(vendorID, file, imageType)
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

	vendorID := mux.Vars(r)["vendorID"]

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

	productID := mux.Vars(r)["id"]

	product.ID, err = strconv.Atoi(productID)
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

	vendorID := mux.Vars(r)["vendorID"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID := mux.Vars(r)["id"]

	err = p.vendorUsecase.DeleteProduct(productID)
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

	vendorID := mux.Vars(r)["vendorID"]

	err = p.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(configs.ImgMaxSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	file, header, err := r.FormFile(configs.ImgFormFileKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	productID := mux.Vars(r)["id"]

	mimeType := header.Header.Get("Content-Type")
	imageType, err := image.CheckMimeType(mimeType)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	imgName, err := p.vendorUsecase.UpdateProductPicture(productID, file, imageType)
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

func (p PartnerDelivery) GetPartnerShops(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	partnerID, ok := r.Context().Value(middleware.UserID(configs.UserID)).(string)
	if !ok {
		err = fmt.Errorf("couldn't get userID from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vendors, err := p.vendorUsecase.GetPartnerShops(partnerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(vendors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
