package handlers

import (
	"ReviewInterfaceAPI/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	DB *sqlx.DB
}

// Structures for JSON API compliance

type RequestBody struct {
	Data models.ReviewData `json:"data"`
}

type ResponseData struct {
	Type       string        `json:"type"`
	ID         int           `json:"id"`
	Attributes models.Review `json:"attributes"`
}

type ResponseBody struct {
	Data ResponseData `json:"data"`
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateReview called")
	productIDStr := chi.URLParam(r, "product_id")
	productId, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	// Decoding
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// checking data type
	if reqBody.Data.Type != "review" {
		http.Error(w, "Invalid data type", http.StatusBadRequest)
		return
	}
	review := reqBody.Data.Attributes
	review.ProductID = productId
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	// Validation
	validate := validator.New()
	if err := validate.Struct(review); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// SQL запрос для вставки нового отзывава
	query := `INSERT INTO reviews (product_id, user_id, rating, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	// SQL запрос для вставки нового отзывава
	var reviewID int
	err = h.DB.QueryRowx(query, review.ProductID, review.UserID, review.Rating, review.Content, review.CreatedAt, review.UpdatedAt).Scan(&reviewID)
	if err != nil {
		log.Printf("Error inserting review into database: %v", err)
		http.Error(w, "Failed to insert review into database", http.StatusInternalServerError)
		return
	}
	review.ID = reviewID

	// Формирование запроса
	respBody := ResponseBody{
		Data: ResponseData{
			Type:       "review",
			ID:         review.ID,
			Attributes: review,
		},
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respBody)
}

func (h *Handler) GetReviews(w http.ResponseWriter, r *http.Request) {
	// Extraction product_id from URL
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	// Handling request params
	includeRatings := r.URL.Query().Get("include") == "ratings"
	sortField := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	// Converting page and limit
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}
	offset := (page - 1) * limit
	// Sql querry with params
	baseQuery := ""
	if includeRatings {
		baseQuery = `SELECT reviews.*, rating FROM reviews WHERE reviews.product_id = $1`
	} else {
		baseQuery = `SELECT * FROM reviews WHERE product_id = $1`
	}
	if sortField != "" {
		baseQuery += fmt.Sprintf(" ORDER BY %s", sortField)
	}

	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// Executing query
	var reviews []models.Review
	err = h.DB.Select(&reviews, baseQuery, productID)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	// Query to get the total number of reviews
	var totalReviews int
	countQuery := `SELECT COUNT(*) FROM reviews WHERE product_id = $1`
	err = h.DB.Get(&totalReviews, countQuery, productID)
	if err != nil {
		http.Error(w, "Failed to query total reviews count", http.StatusInternalServerError)
		return
	}
	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalReviews) / float64(limit)))

	// Pagination metadata
	paginationMeta := map[string]int{
		"totalReviews": totalReviews,
		"totalPages":   totalPages,
		"currentPage":  page,
		"limit":        limit,
	}

	// Forming the response
	response := map[string]interface{}{
		"data": reviews,
		"meta": paginationMeta,
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteReviews(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	query := `DELETE FROM reviews WHERE product_id = $1`
	_, err = h.DB.Exec(query, productID)
	if err != nil {
		http.Error(w, "Failed to delete reviews from database", http.StatusInternalServerError)
		return
	}
	// response about successful deleting
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateCommentById(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	reviewID, err := strconv.Atoi(chi.URLParam(r, "review_id"))
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var req models.ReviewUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateData := req.Data.Attributes

	// Validation
	validate := validator.New()
	if err := validate.Struct(updateData); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

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
	err = h.DB.QueryRow(query, userID, rating, content, productId, reviewID).Scan(&updatedReviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No review found with the given product ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update review in the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("Review with ID %d for product %d updated successfully", updatedReviewID, productId))
}
