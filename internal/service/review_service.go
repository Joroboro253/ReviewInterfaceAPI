package service

import (
	"ReviewInterfaceAPI/internal/models"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
)

type ReviewService struct {
	DB *sqlx.DB
}

func NewReviewService(db *sqlx.DB) *ReviewService {
	return &ReviewService{DB: db}
}

func (s *ReviewService) CreateReview(review *models.Review) (int, error) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.Insert("reviews").
		Columns("product_id", "user_id", "rating", "content", "created_at", "updated_at").
		Values(review.ProductID, review.UserID, review.Rating, review.Content, review.CreatedAt, review.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()
	log.Printf("Executing SQL query: %s with args: %v", query, args)
	if err != nil {
		log.Printf("error building insert SQL query: %v", err)
		return 0, fmt.Errorf("error building insert SQL query: %w", err)
	}

	var reviewID int
	err = s.DB.QueryRow(query, args...).Scan(&reviewID)
	if err != nil {
		log.Printf("error executing insert SQL query: %v", err)
		return 0, fmt.Errorf("error executing insert SQL query: %w", err)
	}

	return reviewID, nil
}

func (s *ReviewService) GetReviewsByProductID(productID int, sortField string, page, limit int) ([]models.Review, int, int, error) {
	// Проверка и корректировка параметров пагинации
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Установка значения по умолчанию
	}
	// Получение самих отзывов
	countBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("COUNT(*)").
		From("reviews").
		Where(squirrel.Eq{"product_id": productID})

	// Подсчет общего количества отзывов
	countQuery, countArgs, err := countBuilder.ToSql()
	log.Printf("Count Query: %s, Args: %v", countQuery, countArgs)
	if err != nil {
		log.Printf("error building count SQL query: %v", err)
		return nil, 0, 0, fmt.Errorf("error building count SQL query: %w", err)
	}

	var totalReviews int
	err = s.DB.Get(&totalReviews, countQuery, countArgs...)
	if err != nil {
		log.Printf("error executing count SQL query: %v", err)
		return nil, 0, 0, fmt.Errorf("error executing count SQL query: %w", err)
	}

	reviewBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("*").
		From("reviews").
		Where(squirrel.Eq{"product_id": productID})

	if sortField != "" {
		reviewBuilder = reviewBuilder.OrderBy(sortField)
	}

	query, args, err := reviewBuilder.Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).ToSql()
	if err != nil {
		log.Printf("error building SQL query: %v", err)
		return nil, 0, 0, fmt.Errorf("error building SQL query: %w", err)
	}

	var reviews []models.Review
	err = s.DB.Select(&reviews, query, args...)
	if err != nil {
		log.Printf("error executing SQL query: %v", err)
		return nil, 0, 0, fmt.Errorf("error executing SQL query: %w", err)
	}

	// Расчет общего количества страниц
	totalPages := int(math.Ceil(float64(totalReviews) / float64(limit)))

	// Возвращение результатов
	return reviews, totalReviews, totalPages, nil
}
func (s *ReviewService) UpdateReview(productId, reviewId int, updateData models.ReviewUpdate) (int, error) {
	// Инициализация SQL-строителя запросов
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("reviews")

	// Добавление условий обновления, если они предоставлены
	if updateData.UserID != nil {
		builder = builder.Set("user_id", *updateData.UserID)
	}
	if updateData.Rating != nil {
		builder = builder.Set("rating", *updateData.Rating)
	}
	if updateData.Content != nil {
		builder = builder.Set("content", *updateData.Content)
	}

	// Добавление условий WHERE и RETURNING
	query, args, err := builder.Where(squirrel.Eq{"id": reviewId, "product_id": productId}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	// Выполнение запроса
	var updatedReviewID int
	err = s.DB.QueryRow(query, args...).Scan(&updatedReviewID)
	if err != nil {
		return 0, err
	}

	return updatedReviewID, nil
}

func (s *ReviewService) DeleteReviewsByProductID(productID int) error {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete("").
		From("reviews").
		Where(squirrel.Eq{"product_id": productID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(query, args...)
	return err
}
