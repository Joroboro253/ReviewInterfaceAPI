package main

import (
	"ReviewInterfaceAPI/internal/repository"
	"ReviewInterfaceAPI/tests/testutils"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

func main() {
	// Configuration
	cfg := repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "bestuser",
		DBName:   "ProductReviewsDB",
	}
	var db *sqlx.DB
	log.Println("Работаем")
	envValue := os.Getenv("ENV")
	log.Println("ENV value:", envValue)

	if os.Getenv("ENV") == "development" {
		log.Println("Зашли в иф")
		db = testutils.CreateTestDB(testutils.Config(cfg))
		log.Println("Создали БД")
	} else {
		log.Println("Зашли в езлэ")

		// Подключение к бд
		var err error
		db, err = repository.NewPostgresDB(cfg)
		if err != nil {
			log.Fatalf("failed to initialize db: %v", err)
		}
	}

	defer db.Close()
	log.Println("Я перед обработчиком")
	// Handler initialization
	//reviewHandler := &handlers.Handler{
	//	DB: db,
	//}

	// Routes setting
	log.Println("Перед роутером")

	//r.Post("/products/{product_id}/reviews", reviewHandler.CreateReview)
	// Other routes

	//Проверка подключения
	app := NewApp(db)
	err := app.Start(":3000")
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
