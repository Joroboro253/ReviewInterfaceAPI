package main

import (
	"ReviewInterfaceAPI/internal/repository"
	"log"
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
	// Подключение к бд
	var err error
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer db.Close()
	//Проверка Соединения
	app := NewApp(db)
	err = app.Start(":3000")
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
