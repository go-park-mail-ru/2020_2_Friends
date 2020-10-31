package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/cart"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"

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

	var cartReq models.CartRequest
	err = json.NewDecoder(r.Body).Decode(&cartReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.cartUsecase.Add(userID, cartReq.ProductID, cartReq.VendorID)
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

	productID, ok := r.URL.Query()["product_id"]
	if !ok {
		err = fmt.Errorf("no query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.cartUsecase.Remove(userID, productID[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
