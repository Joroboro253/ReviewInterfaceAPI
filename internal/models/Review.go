package models

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