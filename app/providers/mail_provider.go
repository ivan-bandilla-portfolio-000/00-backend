package providers

import (
	"log"
	"portfolio-backend/services/email"
)

type MailProvider struct {
	EmailService *email.EmailService
}

func NewMailProvider() *MailProvider {
	emailService, err := email.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}
	return &MailProvider{
		EmailService: emailService,
	}
}
