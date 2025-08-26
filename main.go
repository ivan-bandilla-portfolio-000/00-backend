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

	env := getEnvOrDefault("APP_ENV", "local")

	app := providers.NewAppServiceProvider()

	// Start server
	// Prefer the PORT that hosting provides; fallback to SERVER_PORT then 8080.
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = getEnvOrDefault("SERVER_PORT", "8080")
	}
	serverURL := getEnvOrDefault("SERVER_URL", serverPort)
	log.Printf("Server [%s] started on %s", env, serverURL)

	log.Fatal(http.ListenAndServe(":"+serverPort, app.Handler()))
}
