package model

type ReviewUpdate struct {
	Rating  float64 `json:"rating"`
	Content string  `json:"content"`
}
