package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	log "github.com/friends/pkg/logger"
	"github.com/gorilla/mux"
)

type OrderDelivery struct {
	orderUsecase order.Usecase
}

func New(orderUsecase order.Usecase) OrderDelivery {
	return OrderDelivery{
		orderUsecase: orderUsecase,
	}
}

func (o OrderDelivery) AddOrder(w http.ResponseWriter, r *http.Request) {
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

	order := models.OrderRequest{}
	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order.Sanitize()

	order.CreatedAt = time.Now()

	orderID, err := o.orderUsecase.AddOrder(userID, order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := models.IDResponse{
		ID: orderID,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o OrderDelivery) GetOrder(w http.ResponseWriter, r *http.Request) {
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

	orderID := mux.Vars(r)["id"]

	order, err := o.orderUsecase.GetOrder(userID, orderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
