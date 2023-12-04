package service

import (
	"database/sql"
)

type ReviewService struct {
	DB *sql.DB
}
