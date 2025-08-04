package providers

import (
	"log"
	"portfolio-backend/services/email"
)

type MailProvider struct {
	MailService *email.MailService
}

func NewMailProvider() *MailProvider {
	mailRendererService, err := email.NewMailRendererService()
	if err != nil {
		log.Fatalf("Failed to initialize mail renderer service: %v", err)
	}
	mailService := email.NewMailService(mailRendererService)
	return &MailProvider{
		MailService: mailService,
	}
}
