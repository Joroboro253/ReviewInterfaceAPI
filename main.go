package main

import (
	"ReviewInterfaceAPI/application"
	"context"
	"fmt"
)

func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
