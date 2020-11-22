package usecase

import (
	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/order"
	"github.com/friends/internal/pkg/review"
)

type ReviewUsecase struct {
	reviewRepository review.Repository
	orderRepository  order.Repository
}

func New(reviewRepository review.Repository, orderRepository order.Repository) review.Usecase {
	return ReviewUsecase{
		reviewRepository: reviewRepository,
		orderRepository:  orderRepository,
	}
}

func (r ReviewUsecase) AddReview(review models.Review) error {
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
	return r.reviewRepository.GetVendorReviews(vendorID)
}
