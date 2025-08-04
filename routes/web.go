package routes

import (
	"net/http"
	"portfolio-backend/app/controllers"
)

func RegisterWebRoutes(mux *http.ServeMux, emailController *controllers.EmailController, discordController *controllers.DiscordController) {
	mux.HandleFunc("/send-email", emailController.SendEmail)
	mux.HandleFunc("/discord-webhook", discordController.SendWebhook)
	mux.HandleFunc("/preview-email", emailController.PreviewEmail)
}
