package api

import (
	"encoding/json"
	"net/http"

	"github.com/friends/storage"
	"github.com/gorilla/mux"
)

type VendorService struct {
	db storage.VendorMapDB
}

func (vs VendorService) getVendor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	vendor, err := vs.db.GetVendor(id)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(vendor)
	if err != nil {
		return
	}
}
