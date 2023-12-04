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

type SuccessResponse struct {
	Data struct {
		Type       string      `json:"type"`
		ID         int         `json:"id"`
		Attributes interface{} `json:"attributes,omitempty"`
	} `json:"data"`
}

func sendApiError(w http.ResponseWriter, apiErr *models.APIError) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(apiErr.Status)
	json.NewEncoder(w).Encode(map[string][]models.APIError{"errors": {*apiErr}})
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateReview called")
	productIDStr := chi.URLParam(r, "product_id")
	productId, err := strconv.Atoi(productIDStr)
	if err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}
	// Decoding
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}
	// checking data type
	if reqBody.Data.Type != "review" {
		sendApiError(w, models.ErrInvalidInput)
		return
	}
	review := reqBody.Data.Attributes
	review.ProductID = productId
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	// Validation
	validate := validator.New()
	if err := validate.Struct(review); err != nil {
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}

	// SQL запрос для вставки нового отзывава
	query := `INSERT INTO reviews (product_id, user_id, rating, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	// SQL запрос для вставки нового отзывава
	var reviewID int
	err = h.DB.QueryRowx(query, review.ProductID, review.UserID, review.Rating, review.Content, review.CreatedAt, review.UpdatedAt).Scan(&reviewID)
	if err != nil {
		log.Printf("Error inserting review into database: %v", err)
		sendApiError(w, models.ErrDatabaseProblem)
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
		sendApiError(w, models.ErrInvalidInput)
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
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}
	// Query to get the total number of reviews
	var totalReviews int
	countQuery := `SELECT COUNT(*) FROM reviews WHERE product_id = $1`
	err = h.DB.Get(&totalReviews, countQuery, productID)
	if err != nil {
		sendApiError(w, models.ErrDatabaseProblem)
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
		sendApiError(w, models.ErrInvalidInput)
		return
	}
	query := `DELETE FROM reviews WHERE product_id = $1`
	_, err = h.DB.Exec(query, productID)
	if err != nil {
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}
	// response about successful deleting
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateCommentById(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}

	reviewID, err := strconv.Atoi(chi.URLParam(r, "review_id"))
	if err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}

	var req models.ReviewUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}

	updateData := req.Data.Attributes

	// Validation
	validate := validator.New()
	if err := validate.Struct(updateData); err != nil {
		sendApiError(w, models.ErrInvalidInput)
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
			sendApiError(w, models.ErrReviewNotFound)
			return
		}
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}

	successResp := SuccessResponse{}
	successResp.Data.Type = "review"
	successResp.Data.ID = updatedReviewID
	successResp.Data.Attributes = map[string]interface{}{
		"message": fmt.Sprintf("Review with ID %d for product %d updated successfully", updatedReviewID, productId),
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResp)
}
