package testutils

import "github.com/jmoiron/sqlx"

func SeedTestData(db *sqlx.DB) {
	// Insertion of test data
	db.MustExec(`INSERT INTO reviews (product_id, user_id, rating, content, created_at, updated_at) VALUES (1, 1, 5, 'Perfect product', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`)
}
