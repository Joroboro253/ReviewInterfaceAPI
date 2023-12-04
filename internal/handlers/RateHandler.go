package handlers

// ReviewRatingRequest Review request body
type ReviewRatingRequest struct {
	Rating int `json:"rating" validate:"required,gte=1,lte=5"`
}

// ReviewRatingResponse Response body after review evaluation
type ReviewRatingResponse struct {
	Message string `json:"message"`
}
