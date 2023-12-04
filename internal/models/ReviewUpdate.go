package models

type ReviewUpdate struct {
	UserID  *int     `json:"user_id"`
	Rating  *float64 `json:"rating,omitempty"`
	Content *string  `json:"content,omitempty"`
}
