package routes

import (
	"net/http"
	api_controllers "portfolio-backend/app/controllers/api"
	"portfolio-backend/app/middlewares"
	"portfolio-backend/services/redis"
	"time"
)

func RegisterWebRoutes(mux *http.ServeMux, emailController *api_controllers.EmailController, discordController *api_controllers.DiscordController) {
	rateLimiter := redis.NewUpstashService()
	// Create the middleware (e.g., 5 requests per second, burst 10)
	withRateLimit := func(baseKey string, rps, burst int, ttl time.Duration, handler http.HandlerFunc) http.Handler {
		return middlewares.RateLimitMiddlewareWithKey(rateLimiter, baseKey, rps, burst, ttl)(handler)
	}

	oneHour := time.Hour

	// Wrap your handlers with the middleware
	mux.Handle("/send-email", withRateLimit("send_email", 1, 1, oneHour, emailController.SendEmail))
	mux.Handle("/discord-webhook", withRateLimit("discord_webhook", 1, 1, oneHour, discordController.SendWebhook))
	mux.HandleFunc("/preview-email", emailController.PreviewEmail)
}
