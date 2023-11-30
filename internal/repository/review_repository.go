package repository

import (
	"ReviewInterfaceAPI/internal/models"
	"context"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review models.Review) (int, error)
	GetReview(ctx context.Context, id int) (*models.Review, error)
	UpdateReview(ctx context.Context, id int, reviewUpdate models.ReviewUpdate) error
	DeleteReview(ctx context.Context, id int) error
	ListReviews(ctx context.Context, page int, size int) ([]models.Review, error)
}
