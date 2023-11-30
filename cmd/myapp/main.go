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
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}
	defer db.Close()

	//Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	} else {
		log.Println("Successfully connected to DB")
	}

}
