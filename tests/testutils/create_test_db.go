package testutils

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func CreateTestDB(cfg Config) *sqlx.DB {
	// Create db sqlite
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password)

	db, _ := sqlx.Open("postgres", psqlInfo)

	//Creation a product table
	//productSchema := `
	//CREATE TABLE products (
	//  id SERIAL PRIMARY KEY,
	//  name TEXT,
	//  description TEXT
	//);`
	//db.MustExec(productSchema)
	//
	//Creation review table
	//reviewSchema := `
	//CREATE TABLE reviews (
	//  id SERIAL PRIMARY KEY,
	//  product_id INTEGER,
	//  user_id INTEGER,
	//  rating INTEGER,
	//  content TEXT,
	//	created_at TIMESTAMP,
	//	updated_at TIMESTAMP
	//);`
	//db.MustExec(reviewSchema)

	// Insertion of test prod
	//db.MustExec(`INSERT INTO products (name, description) VALUES ('Тестовый продукт 1', 'Описание тестового продукта 1')`)
	//db.MustExec(`INSERT INTO products (name, description) VALUES ('Тестовый продукт 2', 'Описание тестового продукта 2')`)

	return db
}
