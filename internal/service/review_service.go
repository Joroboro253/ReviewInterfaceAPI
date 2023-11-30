package service

import (
	"ReviewInterfaceAPI/internal/model"
	"ReviewInterfaceAPI/internal/repository"
	"context"
)

// for handling business logic related to rewies
type ReviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		repo: repo,
	}
}

// Create a new review
func (s *ReviewService) CreateReview(ctx context.Context, review model.Review) (int, error) {
	return s.repo.CreateReview(ctx, review)
}

// Get rewiew by ID
func (s *ReviewService) GetReview(ctx context.Context, id int) (*model.Review, error) {
	return s.repo.GetReview(ctx, id)
}

// Update review by ID
func (s *ReviewService) UpdateReview(ctx context.Context, id int, reviewUpdate model.ReviewUpdate) error {
	return s.repo.UpdateReview(ctx, id, reviewUpdate)
}

func (s *ReviewService) DeleteReview(ctx context.Context, id int) error {
	return s.repo.DeleteReview(ctx, id)
}

func (s *ReviewService) ListReviews(ctx context.Context, page int, size int) ([]model.Review, error) {
	return s.repo.ListReviews(ctx, page, size)
}
