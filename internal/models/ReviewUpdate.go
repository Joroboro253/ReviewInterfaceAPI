package models

type ReviewUpdate struct {
	Rating  *float64 `json:"rating,omitempty"`
	Content *string  `json:"content,omitempty"`
}
