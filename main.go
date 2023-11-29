package main

import (
	"ReviewInterfaceAPI/application"
	"context"
	"fmt"
	"log"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("failed to start app:", err)
	}

	//db, err := repository.NewPostgresDB(repository.Config{
	//	Host:     "localhost",
	//	Port:     "5436",
	//	Username: "postgres",
	//	Password: "bestuser",
	//	DBName:   "ProductReviewsDB",
	//	SSLMode:  "disable",
	//})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

}
