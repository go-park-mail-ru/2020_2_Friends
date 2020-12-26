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
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/vendors"
	websocketpool "github.com/friends/internal/pkg/websocketPool"
	ownErr "github.com/friends/pkg/error"
	log "github.com/friends/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type OrderDelivery struct {
	orderUsecase  order.Usecase
	vendorUsecase vendors.Usecase
	websocketPool websocketpool.WebsocketPool
}

func New(
	orderUsecase order.Usecase, vendorUsecase vendors.Usecase, websocketPool websocketpool.WebsocketPool,
) OrderDelivery {
	return OrderDelivery{
		orderUsecase:  orderUsecase,
		vendorUsecase: vendorUsecase,
		websocketPool: websocketPool,
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
	log.DataLog(r, order)

	order.CreatedAt = time.Now()

	orderID, err := o.orderUsecase.AddOrder(userID, order)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusConflict)
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

	orderID, ok := mux.Vars(r)["id"]
	if !ok {
		err = fmt.Errorf("no id in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := o.orderUsecase.GetOrder(userID, orderID)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o OrderDelivery) GetUserOrders(w http.ResponseWriter, r *http.Request) {
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

	orders, err := o.orderUsecase.GetUserOrders(userID)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o OrderDelivery) GetVendorOrders(w http.ResponseWriter, r *http.Request) {
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

	vendorID, ok := mux.Vars(r)["id"]
	if !ok {
		err = fmt.Errorf("no vendor id in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.vendorUsecase.CheckVendorOwner(partnerID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := o.orderUsecase.GetVendorOrders(vendorID)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o OrderDelivery) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
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

	vendorID, ok := mux.Vars(r)["vendorID"]
	if !ok {
		err = fmt.Errorf("no vendor id in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.vendorUsecase.CheckVendorOwner(partnerID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := models.OrderStatusRequest{}
	err = json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	status.Sanitize()
	log.DataLog(r, status)

	orderID, ok := mux.Vars(r)["id"]
	if !ok {
		err = fmt.Errorf("no order id in path")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderIDInt, err := strconv.Atoi(orderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.orderUsecase.UpdateOrderStatus(orderID, status.Status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	clientID, err := o.orderUsecase.GetUserIDFromOrder(orderIDInt)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	wsConn, ok := o.websocketPool.Get(clientID)
	if !ok {
		return
	}

	vendor, err := o.vendorUsecase.GetVendorInfo(vendorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orderMessage := models.OrderStatusMessage{
		Type:          "status",
		OrderID:       orderID,
		VendorName:    vendor.Name,
		VendorPicture: vendor.Picture,
		Status:        status.Status,
	}

	msgJSON, err := json.Marshal(orderMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = wsConn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
