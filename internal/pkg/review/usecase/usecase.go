package usecase

import (
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/review"
	"github.com/friends/internal/pkg/vendors"
	ownErr "github.com/friends/pkg/error"
)

type ReviewUsecase struct {
	reviewRepository  review.Repository
	orderRepository   order.Repository
	profileRepository profile.Repository
	vendorRepository  vendors.Repository
}

func New(
	reviewRepository review.Repository, orderRepository order.Repository,
	profileRepository profile.Repository, vendorRepository vendors.Repository,
) review.Usecase {
	return ReviewUsecase{
		reviewRepository:  reviewRepository,
		orderRepository:   orderRepository,
		profileRepository: profileRepository,
		vendorRepository:  vendorRepository,
	}
}

func (r ReviewUsecase) AddReview(review models.Review) error {
	order, err := r.orderRepository.GetOrder(strconv.Itoa(review.OrderID))
	if err != nil {
		return err
	}

	if order.Reviewed {
		return ownErr.NewClientError(fmt.Errorf("review already exists for order with id = %v", order.ID))
	}

	if review.UserID != strconv.Itoa(order.UserID) {
		return ownErr.NewClientError(fmt.Errorf("the order does not belong to the user"))
	}

	review.VendorID = order.VendorID

	err = r.reviewRepository.AddReview(review)
	if err != nil {
		return err
	}

	err = r.orderRepository.SetOrderReviewStatus(review.OrderID, true)
	if err != nil {
		return err
	}

	return nil
}

func (r ReviewUsecase) GetUserReviews(userID string) ([]models.Review, error) {
	return r.reviewRepository.GetUserReviews(userID)
}

func (r ReviewUsecase) GetVendorReviews(vendorID string) (models.VendorReviewsResponse, error) {
	idInt, err := strconv.Atoi(vendorID)
	if err != nil {
		return models.VendorReviewsResponse{}, err
	}
	vendor, err := r.vendorRepository.Get(idInt)
	if err != nil {
		return models.VendorReviewsResponse{}, fmt.Errorf("couldn't get vendor: %w", err)
	}

	reviews, err := r.reviewRepository.GetVendorReviews(vendorID)

	if err != nil {
		return models.VendorReviewsResponse{}, err
	}

	for idx, review := range reviews {
		name, err := r.profileRepository.GetUsername(review.UserID)
		if err != nil {
			return models.VendorReviewsResponse{}, err
		}

		reviews[idx].Username = name
	}

	resp := models.VendorReviewsResponse{
		VendorName:    vendor.Name,
		VendorPicture: vendor.Picture,
		Reviews:       reviews,
	}
	return resp, nil
}
