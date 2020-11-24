package review

import "github.com/friends/internal/pkg/models"

type Repository interface {
	AddReview(models.Review) error
	GetUserReviews(userID string) ([]models.Review, error)
	GetVendorReviews(vendorID string) ([]models.Review, error)
}
