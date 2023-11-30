package repository

import (
	"ReviewInterfaceAPI/internal/model"
	"context"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review model.Review) (int, error)
	GetReview(ctx context.Context, id int) (*model.Review, error)
	UpdateReview(ctx context.Context, id int, reviewUpdate model.ReviewUpdate) error
	DeleteReview(ctx context.Context, id int) error
	ListReviews(ctx context.Context, page int, size int) ([]model.Review, error)
}
