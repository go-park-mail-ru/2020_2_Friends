package repository

import (
	"database/sql"
	"fmt"

	"github.com/friends/internal/pkg/models"
	"github.com/friends/internal/pkg/review"
)

type ReviewRepository struct {
	db *sql.DB
}

func New(db *sql.DB) review.Repository {
	return ReviewRepository{
		db: db,
	}
}

func (r ReviewRepository) AddReview(review models.Review) error {
	_, err := r.db.Exec(
		`INSERT INTO reviews (userID, orderID, vendorID, rating, review_text, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		review.UserID, review.OrderID, review.VendorID, review.Rating, review.Text, review.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("couldn't insert review: %w", err)
	}

	return nil
}

func (r ReviewRepository) GetUserReviews(userID string) ([]models.Review, error) {
	rows, err := r.db.Query(
		`SELECT userID, orderID, rating, review_text, created_at FROM reviews
		WHERE userID = $1`,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get user reviews: %w", err)
	}
	defer rows.Close()

	reviews := make([]models.Review, 0)
	for rows.Next() {
		review := models.Review{}
		err = rows.Scan(&review.UserID, &review.OrderID, &review.Rating, &review.Text, &review.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("couldn't get user review: %w", err)
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r ReviewRepository) GetVendorReviews(vendorID string) ([]models.Review, error) {
	rows, err := r.db.Query(
		`SELECT userID, orderID, rating, review_text, created_at FROM reviews
		WHERE vendorID = $1`,
		vendorID,
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't get vendor reviews: %w", err)
	}
	defer rows.Close()

	reviews := make([]models.Review, 0)
	for rows.Next() {
		review := models.Review{}
		err = rows.Scan(&review.UserID, &review.OrderID, &review.Rating, &review.Text, &review.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("couldn't get vendor review: %w", err)
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}
