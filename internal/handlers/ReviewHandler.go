package handlers

import (
	"ReviewInterfaceAPI/internal/models"
	"ReviewInterfaceAPI/internal/service"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"log"
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

	reviewService := service.NewReviewService(h.DB)
	reviewID, err := reviewService.CreateReview(&review)
	if err != nil {
		log.Printf("Error inserting review into database: %v", err)
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}
	review.ID = reviewID

	// Query generation
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
	// Extracting product_id from URL
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}
	log.Printf("Product id: %v", productID)
	// Query parameter processing
	sortField := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Conversion page and limit
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // default value
	}

	reviewService := service.NewReviewService(h.DB)
	reviews, totalReviews, totalPages, err := reviewService.GetReviewsByProductID(productID, sortField, page, limit)
	if err != nil {
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}

	// Pagination metadata
	paginationMeta := map[string]int{
		"totalReviews": totalReviews,
		"totalPages":   totalPages,
		"currentPage":  page,
		"limit":        limit,
	}

	// Response formation
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
	reviewService := service.NewReviewService(h.DB)
	err = reviewService.DeleteReviewsByProductID(productID)
	if err != nil {
		sendApiError(w, models.ErrDatabaseProblem)
		return
	}
	// response about successful deleting
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateReviewById(w http.ResponseWriter, r *http.Request) {
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
	// Decoding
	var req models.ReviewUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}
	// Validation
	updateData := req.Data.Attributes
	validate := validator.New()
	if err := validate.Struct(updateData); err != nil {
		sendApiError(w, models.ErrInvalidInput)
		return
	}

	reviewService := service.NewReviewService(h.DB)
	updatedReviewID, err := reviewService.UpdateReview(productId, reviewID, updateData)
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
