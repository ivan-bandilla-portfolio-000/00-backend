package main

import (
	"log"
	"net/http"
	"os"

	"portfolio-backend/app/providers"

	"github.com/joho/godotenv"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	app := providers.NewAppServiceProvider()

	// Start server
	port := getEnvOrDefault("PORT", "8080")
	log.Printf("Server started on :%s", port)

	if app.MailProvider.EmailService != nil {
		log.Println("Email templates loaded successfully")
	} else {
		log.Println("Email templates will use fallback HTML")
	}

	log.Fatal(http.ListenAndServe(":"+port, app.Handler()))
}
