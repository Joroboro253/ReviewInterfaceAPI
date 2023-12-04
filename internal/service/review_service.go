package service

import (
	"ReviewInterfaceAPI/internal/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"
)

type ReviewService struct {
	DB *sqlx.DB
}

func NewReviewService(db *sqlx.DB) *ReviewService {
	return &ReviewService{DB: db}
}

func (s *ReviewService) CreateReview(review *models.Review) (int, error) {
	query := `INSERT INTO reviews (product_id, user_id, rating, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var reviewID int
	err := s.DB.QueryRow(query, review.ProductID, review.UserID, review.Rating, review.Content, review.CreatedAt, review.UpdatedAt).Scan(&reviewID)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}

func (s *ReviewService) GetReviewsByProductID(productID int, includeRatings bool, sortField string, page, limit int) ([]models.Review, int, int, error) {
	//Construct base query
	baseQuery := ""
	if includeRatings {
		baseQuery = `SELECT reviews.*, rating FROM reviews WHERE reviews.product_id = $1`
	} else {
		baseQuery = `SELECT * FROM reviews WHERE product_id = $1`
	}
	if sortField != "" {
		baseQuery += fmt.Sprintf(" ORDER BY %s", sortField)
	}
	offset := (page - 1) * limit
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Execute query
	var reviews []models.Review
	err := s.DB.Select(&reviews, baseQuery, productID)
	if err != nil {
		return nil, 0, 0, err
	}
	// Query to get the total number of reviews
	var totalReviews int
	countQuery := `SELECT COUNT(*) FROM reviews WHERE product_id = $1`
	err = s.DB.Get(&totalReviews, countQuery, productID)
	if err != nil {
		return nil, 0, 0, err
	}
	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalReviews) / float64(limit)))

	return reviews, totalReviews, totalPages, nil
}

func (s *ReviewService) UpdateReview(productId, reviewId int, updateData models.ReviewUpdate) (int, error) {
	rating := sql.NullInt64{Valid: updateData.Rating != nil}
	if rating.Valid {
		rating.Int64 = int64(*updateData.Rating)
	}
	content := sql.NullString{Valid: updateData.Content != nil}
	if content.Valid {
		content.String = *updateData.Content
	}
	userID := sql.NullInt64{Valid: updateData.UserID != nil}
	if userID.Valid {
		userID.Int64 = int64(*updateData.UserID)
	}
	// sql query for updating
	query := `UPDATE reviews SET user_id = COALESCE($1, user_id), rating = COALESCE($2, rating), content = COALESCE($3, content) WHERE product_id = $4 AND id = $5 RETURNING id`
	var updatedReviewID int
	err := s.DB.QueryRow(query, userID, rating, content, productId, reviewId).Scan(&updatedReviewID)
	if err != nil {
		return 0, err
	}
	return updatedReviewID, nil
}

func (s *ReviewService) DeleteReviewsByProductID(productID int) error {
	query := `DELETE FROM reviews WHERE product_id = $1`
	_, err := s.DB.Exec(query, productID)
	return err
}
