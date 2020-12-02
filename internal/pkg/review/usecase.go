package review

import "github.com/friends/internal/pkg/models"

//go:generate mockgen -destination=./usecase_mock.go -package=review github.com/friends/internal/pkg/review Usecase
type Usecase interface {
	AddReview(models.Review) error
	GetUserReviews(userID string) ([]models.Review, error)
	GetVendorReviews(vendorID string) (models.VendorReviewsResponse, error)
}
