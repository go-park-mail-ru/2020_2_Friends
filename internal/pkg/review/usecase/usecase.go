package usecase

import (
	"fmt"
	"strconv"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/profile"
	"github.com/friends/internal/pkg/review"

	ownErr "github.com/friends/pkg/error"
)

type ReviewUsecase struct {
	reviewRepository  review.Repository
	orderRepository   order.Repository
	profileRepository profile.Repository
}

func New(
	reviewRepository review.Repository, orderRepository order.Repository, profileRepository profile.Repository,
) review.Usecase {
	return ReviewUsecase{
		reviewRepository:  reviewRepository,
		orderRepository:   orderRepository,
		profileRepository: profileRepository,
	}
}

func (r ReviewUsecase) AddReview(review models.Review) error {
	isUserOrder := r.orderRepository.CheckOrderByUser(review.UserID, strconv.Itoa(review.OrderID))
	if !isUserOrder {
		return ownErr.NewClientError(fmt.Errorf("the order does not belong to the user"))
	}

	vendorID, err := r.orderRepository.GetVendorIDFromOrder(review.OrderID)
	if err != nil {
		return err
	}

	review.VendorID = vendorID

	return r.reviewRepository.AddReview(review)
}

func (r ReviewUsecase) GetUserReviews(userID string) ([]models.Review, error) {
	return r.reviewRepository.GetUserReviews(userID)
}

func (r ReviewUsecase) GetVendorReviews(vendorID string) ([]models.Review, error) {
	reviews, err := r.reviewRepository.GetVendorReviews(vendorID)

	if err != nil {
		return nil, err
	}

	for idx, review := range reviews {
		name, err := r.profileRepository.GetUsername(review.UserID)
		if err != nil {
			return nil, fmt.Errorf("couldn't get username: %w", err)
		}

		reviews[idx].Username = name
	}

	return reviews, nil
}
