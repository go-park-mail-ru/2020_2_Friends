package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/chat"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/vendors"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	ownErr "github.com/friends/pkg/error"
	log "github.com/friends/pkg/logger"
)

type ChatDelivery struct {
	chatUsecase   chat.Usecase
	orderUsecase  order.Usecase
	vendorUsecase vendors.Usecase
	upgrader      websocket.Upgrader
	socketPool    map[string]*websocket.Conn
	mux           *sync.RWMutex
}

func New(chatUsecase chat.Usecase, orderUsecase order.Usecase, vendorUsecase vendors.Usecase) ChatDelivery {
	return ChatDelivery{
		chatUsecase:   chatUsecase,
		orderUsecase:  orderUsecase,
		vendorUsecase: vendorUsecase,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		socketPool: make(map[string]*websocket.Conn),
		mux:        &sync.RWMutex{},
	}
}

func (c ChatDelivery) Upgrade(w http.ResponseWriter, r *http.Request) {
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

	ws, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	c.addConn(ws, userID)

	c.read(ws, userID)

	c.deleteConn(userID)
}

func (c ChatDelivery) addConn(ws *websocket.Conn, userID string) {
	c.mux.RLock()
	c.socketPool[userID] = ws
	c.mux.RUnlock()
}

func (c ChatDelivery) deleteConn(userID string) {
	c.mux.RLock()
	delete(c.socketPool, userID)
	c.mux.RUnlock()
}

func (c ChatDelivery) read(ws *websocket.Conn, userID string) {
	for {
		_, msgJSON, err := ws.ReadMessage()
		if err != nil {
			continue
		}

		msg := models.Message{}
		err = json.Unmarshal(msgJSON, &msg)
		if err != nil {
			continue
		}
		msg.UserID = userID
		msg.SentAt = time.Now()

		err = c.chatUsecase.Save(msg)
		if err != nil {
			return
		}

		customerID, err := c.orderUsecase.GetUserIDFromOrder(msg.OrderID)
		if err != nil {
			continue
		}

		vendorID, err := c.orderUsecase.GetVendorIDFromOrder(msg.OrderID)
		if err != nil {
			continue
		}

		partnerID, err := c.vendorUsecase.GetVendorOwner(vendorID)
		if err != nil {
			continue
		}

		if userID == customerID {
			c.write(partnerID, msgJSON)
		} else if userID == partnerID {
			c.write(customerID, msgJSON)
		}
	}
}

func (c ChatDelivery) write(userID string, text []byte) {
	conn, ok := c.socketPool[userID]
	if !ok {
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, text)
	if err != nil {
		fmt.Println(err)
	}
}

func (c ChatDelivery) GetChat(w http.ResponseWriter, r *http.Request) {
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

	orderIDStr, ok := mux.Vars(r)["id"]
	if !ok {
		err = fmt.Errorf("no id in url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userIDFromDB, err := c.orderUsecase.GetUserIDFromOrder(orderID)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}

	if userID != userIDFromDB {
		vendorID, err := c.orderUsecase.GetVendorIDFromOrder(orderID)
		if err != nil {
			ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
			return
		}

		err = c.vendorUsecase.CheckVendorOwner(userID, strconv.Itoa(vendorID))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	msgs, err := c.chatUsecase.GetChat(orderID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(msgs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c ChatDelivery) GetVendorChats(w http.ResponseWriter, r *http.Request) {
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

	vendorID, ok := mux.Vars(r)["id"]
	if !ok {
		err = fmt.Errorf("no id in url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.vendorUsecase.CheckVendorOwner(userID, vendorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chats, err := c.chatUsecase.GetVendorChats(vendorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
