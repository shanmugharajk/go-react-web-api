package main

import (
	"log"

	"github.com/shanmugharajk/go-react-web-api/api/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("Application error:", err)
	}
}
