package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/friends/internal/pkg/vendors"
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
	strID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vendor, err := v.vendorUsecase.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(vendor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
