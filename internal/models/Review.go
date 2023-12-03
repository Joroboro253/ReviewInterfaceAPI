package models

import "time"

// structure of entity
type Review struct {
	ID        int       `json:"id" db:"id"`
	ProductID int       `json:"product_id" db:"product_id" validate:"required,gte=1"`
	UserID    int       `json:"user_id" db:"user_id" validate:"required,gte=1"`
	Rating    *int      `json:"rating" db:"rating" validate:"omitempty,gte=1,lte=5"`
	Content   string    `json:"content" db:"content" validate:"required,max=1000"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
