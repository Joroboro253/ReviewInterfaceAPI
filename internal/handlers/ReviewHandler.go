package handlers

import (
	"ReviewInterfaceAPI/internal/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Repo *sqlx.DB
}

// структура для использования в обработчике
type ReviewRequest struct {
	ProductID int    `json:"productId"`
	UserID    int    `json:"userID"`
	Rating    int    `json:"rating"`
	Content   string `json:"content"`
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	// SQL запрос для вставки нового отзывава
	query := `INSERT INTO reviews (product_id, user_id, rating, content, created_at, updated_at) VALUES (:product_id, :user_id, :rating, :content, :created_at, :updated_at)`

	// Data insertion
	_, err := h.Repo.NamedExec(query, review)
	if err != nil {
		http.Error(w, "Failed to insert review into database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func (h *Handler) GetReviews(w http.ResponseWriter, r *http.Request) {
	// Extraction product_id from URL
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var reviews []models.Review
	// Sql querry for gettting reviews with help product_id
	query := `SELECT * FROM reviews WHERE product_id = $1`
	err = h.Repo.Select(&reviews, query, productID)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) DeleteReviews(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	query := `DELETE FROM reviews WHERE product_id = $1`
	_, err = h.Repo.Exec(query, productID)
	if err != nil {
		http.Error(w, "Failed to delete reviews from database", http.StatusInternalServerError)
		return
	}
	// response about successful deleting
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateReviews(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updateData models.ReviewUpdate
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sql query for updating
	if updateData.Rating != nil {
		query := `UPDATE reviews SET rating = $1 WHERE product_id = $2`
		_, err = h.Repo.Exec(query, *updateData.Rating, productID)
		if err != nil {
			http.Error(w, "Failed to update reviews in the database", http.StatusInternalServerError)
			return
		}
	}
	if updateData.Content != nil {
		query := `UPDATE reviews SET content = $q WHERE product_id = $2`
		_, err = h.Repo.Exec(query, *updateData.Content, productID)
		if err != nil {
			http.Error(w, "Failed to update reviews in the database", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Reviews updated successfully")
}
