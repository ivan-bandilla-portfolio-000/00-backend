package main

import (
	"log"
	"net/http"
	"os"

	"portfolio-backend/app/providers"

	"portfolio-backend/bootstrap"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	bootstrap.LoadEnv()

	app := providers.NewAppServiceProvider()

	// Start server
	port := getEnvOrDefault("PORT", "8080")
	log.Printf("Server started on :%s", port)

	log.Fatal(http.ListenAndServe(":"+port, app.Handler()))
}
