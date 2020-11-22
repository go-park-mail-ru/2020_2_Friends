package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/review"
	"github.com/friends/internal/pkg/vendors"
	"github.com/gorilla/mux"

	ownErr "github.com/friends/pkg/error"
	log "github.com/friends/pkg/logger"
)

type ReviewDelivery struct {
	reviewUsecase review.Usecase
	vendorUsecase vendors.Usecase
}

func New(reviewUsecase review.Usecase, vendorUsecase vendors.Usecase) ReviewDelivery {
	return ReviewDelivery{
		reviewUsecase: reviewUsecase,
		vendorUsecase: vendorUsecase,
	}
}

func (rd ReviewDelivery) AddReview(w http.ResponseWriter, r *http.Request) {
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

	review := models.Review{}
	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	review.Sanitize()

	review.UserID = userID
	review.CreatedAt = time.Now()

	err = rd.reviewUsecase.AddReview(review)
	if err != nil {
		ownErr.HandleErrorAndWriteResponse(w, err)
		return
	}
}

func (re ReviewDelivery) GetUserReviews(w http.ResponseWriter, r *http.Request) {
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

	reviews, err := re.reviewUsecase.GetUserReviews(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (re ReviewDelivery) GetVendorReviews(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.ErrorLogWithCtx(r.Context(), err)
		}
	}()

	vendorID, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reviews, err := re.reviewUsecase.GetVendorReviews(vendorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
