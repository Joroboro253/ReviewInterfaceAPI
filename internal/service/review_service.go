package service

import (
	"database/sql"
)

// for handling business logic related to rewies
type ReviewService struct {
	DB *sql.DB
}

// Create a new review
//func (s *ReviewService) CreateReview(productId int, reviewData models.ReviewData) (models.Review, error) {
//	review := models.Review{
//		ProductID: productId,
//		UserID:    reviewData.UserID,
//		Rating:    reviewData.Rating,
//		Content:   reviewData.Content,
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//	}
//	query := `INSERT INTO reviews (product_id, user_id, rating, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
//	err := s.DB.QueryRowx(query, review.ProductID, review.UserID, review.Rating, review.Content, review.CreatedAt, review.UpdatedAt).Scan(&reviewID)
//	if err != nil {
//		log.Printf("Error inserting review into database: %v", err)
//		return Review{}, fmt.Errorf("failed to insert review into database: %v", err)
//	}
//
//	return review, nil
//}

//// Get rewiew by ID
//func (s *ReviewService) GetReview(ctx context.Context, id int) (*models.Review, error) {
//	return s.repo.GetReview(ctx, id)
//}
//
//// Update review by ID
//func (s *ReviewService) UpdateReview(ctx context.Context, id int, reviewUpdate models.ReviewUpdate) error {
//	return s.repo.UpdateReview(ctx, id, reviewUpdate)
//}
//
//func (s *ReviewService) DeleteReview(ctx context.Context, id int) error {
//	return s.repo.DeleteReview(ctx, id)
//}
//
//func (s *ReviewService) ListReviews(ctx context.Context, page int, size int) ([]models.Review, error) {
//	return s.repo.ListReviews(ctx, page, size)
//}
