package review

import "github.com/friends/internal/pkg/models"

type Usecase interface {
	AddReview(models.Review) error
	GetUserReviews(userID string) ([]models.Review, error)
	GetVendorReviews(vendorID string) (models.VendorReviewsResponse, error)
}
