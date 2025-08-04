package providers

import (
	"net/http"
	api_controllers "portfolio-backend/app/controllers/api"
	"portfolio-backend/app/middlewares"
	"portfolio-backend/config"
	"portfolio-backend/routes"
)

type AppServiceProvider struct {
	Mux               *http.ServeMux
	CorsProvider      *CorsProvider
	MailProvider      *MailProvider
	EmailController   *api_controllers.EmailController
	DiscordController *api_controllers.DiscordController
}

func NewAppServiceProvider() *AppServiceProvider {
	corsConfig := config.LoadCORSConfig()
	corsProvider := NewCorsProvider(corsConfig)
	mux := http.NewServeMux()
	mailProvider := NewMailProvider()

	emailController := api_controllers.NewEmailController(mailProvider.MailService)
	discordController := api_controllers.NewDiscordController()

	routes.RegisterAuthRoutes(mux)
	routes.RegisterWebRoutes(mux, emailController, discordController)

	return &AppServiceProvider{
		Mux:               mux,
		CorsProvider:      corsProvider,
		MailProvider:      mailProvider,
		EmailController:   emailController,
		DiscordController: discordController,
	}
}

func (asp *AppServiceProvider) Handler() http.Handler {
	// Wrap the mux with the global rate limiter, then adapt to HandlerFunc for CORS
	return asp.CorsProvider.Handler(
		func(w http.ResponseWriter, r *http.Request) {
			middlewares.GlobalRateLimiter(asp.Mux).ServeHTTP(w, r)
		},
	)
}
