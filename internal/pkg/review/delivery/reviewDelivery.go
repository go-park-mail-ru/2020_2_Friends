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
	ownErr "github.com/friends/pkg/error"
	log "github.com/friends/pkg/logger"
	"github.com/gorilla/mux"
)

type ReviewDelivery struct {
	reviewUsecase review.Usecase
}

func New(reviewUsecase review.Usecase) ReviewDelivery {
	return ReviewDelivery{
		reviewUsecase: reviewUsecase,
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
		ownErr.HandleErrorAndWriteResponse(w, err, http.StatusBadRequest)
		return
	}
}

func (rd ReviewDelivery) GetUserReviews(w http.ResponseWriter, r *http.Request) {
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

	reviews, err := rd.reviewUsecase.GetUserReviews(userID)
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

func (rd ReviewDelivery) GetVendorReviews(w http.ResponseWriter, r *http.Request) {
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

	reviews, err := rd.reviewUsecase.GetVendorReviews(vendorID)
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
