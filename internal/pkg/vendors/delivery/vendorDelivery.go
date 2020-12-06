package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/vendors"
	log "github.com/friends/pkg/logger"
	"github.com/gorilla/mux"
)

type VendorDelivery struct {
	vendorUsecase vendors.Usecase
}

func NewVendorDelivery(usecase vendors.Usecase) VendorDelivery {
	return VendorDelivery{
		vendorUsecase: usecase,
	}
}

func (v VendorDelivery) GetVendor(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	strID, ok := mux.Vars(r)["id"]
	if !ok {
		err = fmt.Errorf("no id in url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vendor, err := v.vendorUsecase.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(vendor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (v VendorDelivery) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	vendors, err := v.vendorUsecase.GetAll()
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

func (v VendorDelivery) GetNearest(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	longitudeQueryParam, ok := r.URL.Query()[configs.Longitude]
	if !ok {
		err = fmt.Errorf("no query param longitude")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	latitudeQueryParam, ok := r.URL.Query()[configs.Latitude]
	if !ok {
		err = fmt.Errorf("no query param latitude")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeQueryParam[0], 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(latitudeQueryParam[0], 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vendors, err := v.vendorUsecase.GetNearest(longitude, latitude)
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
