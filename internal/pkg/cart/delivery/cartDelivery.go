package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/cart"
	"github.com/friends/internal/pkg/middleware"
	log "github.com/friends/pkg/logger"
)

type CartDelivery struct {
	cartUsecase cart.Usecase
}

func NewCartDelivery(cartUsecase cart.Usecase) CartDelivery {
	return CartDelivery{
		cartUsecase: cartUsecase,
	}
}

func (c CartDelivery) AddToCart(w http.ResponseWriter, r *http.Request) {
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

	productID, ok := r.URL.Query()[configs.ProductID]
	if !ok {
		err = fmt.Errorf("no query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	countStr, ok := r.URL.Query()["count"]
	if !ok {
		err = fmt.Errorf("no query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(countStr[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.cartUsecase.Add(userID, productID[0], count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c CartDelivery) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
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

	productID, ok := r.URL.Query()[configs.ProductID]
	if !ok {
		err = fmt.Errorf("no query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.cartUsecase.Remove(userID, productID[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c CartDelivery) GetCart(w http.ResponseWriter, r *http.Request) {
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

	products, err := c.cartUsecase.Get(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
