package model

import "time"

// structure of entity
type Review struct {
	ID        int       `json:"id"`
	ProductID int       `json:"productId"`
	UserID    int       `json:"userID"`
	Rating    int       `json:"rating"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ReviewUpdate struct {
	Rating  float64 `json:"rating"`
	Content string  `json:"content"`
}
