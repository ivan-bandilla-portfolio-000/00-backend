package email

import (
	"log"
	"os"
	validation_email "portfolio-backend/services/validation/email"
	"portfolio-backend/utils"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	MailRendererService *MailRendererService
}

func NewMailService(mailRendererService *MailRendererService) *MailService {
	return &MailService{MailRendererService: mailRendererService}
}

type ContactRequest struct {
	From    string
	Subject string
	Body    string
	Name    string
}

func (cs *MailService) SendContactEmail(req ContactRequest) error {
	// Validate email
	_, err := validation_email.ValidateEmail(req.From)
	if err != nil {
		return err
	}

	// Render sender info
	var senderInfo strings.Builder
	err = cs.MailRendererService.Templates().ExecuteTemplate(&senderInfo, "components/ui/sender.tmpl", map[string]interface{}{
		"FromEmail": req.From,
		"Subject":   req.Subject,
	})
	if err != nil {
		log.Printf("Failed to render sender info: %v", err)
	}

	emailBodyWithSenderInfo := senderInfo.String() + req.Body

	emailData := EmailData{
		AppName:     utils.GetEnvOrDefault("APP_NAME", "Portfolio"),
		Subject:     req.Subject,
		SiteURL:     utils.GetEnvOrDefault("APP_URL", "https://yourdomain.com"),
		HeaderTitle: utils.GetEnvOrDefault("APP_NAME", "Portfolio"),
		Year:        time.Now().Year(),
		FromName:    utils.GetNameFromEmail(req.From, req.Name),
		FromEmail:   req.From,
		Body:        emailBodyWithSenderInfo,
	}

	htmlBody, err := cs.MailRendererService.RenderContactEmail(emailData)
	if err != nil {
		log.Printf("Failed to render email template: %v", err)
		htmlBody = generateFallbackHTML(req)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", os.Getenv("THIS_PORTFOLIO_CONTACT_EMAIL"))
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", htmlBody)
	plainBody := utils.StripHTMLTags(htmlBody)
	m.AddAlternative("text/plain", plainBody)

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func (cs *MailService) GeneratePreviewEmail() (string, error) {
	emailData := EmailData{
		AppName:     utils.GetEnvOrDefault("APP_NAME", "Ivan Bandilla Portfolio"),
		Subject:     "Contact Form Preview",
		SiteURL:     utils.GetEnvOrDefault("APP_URL", "https://ivanbandilla.dev"),
		HeaderTitle: utils.GetEnvOrDefault("APP_NAME", "Ivan Bandilla"),
		Year:        time.Now().Year(),
		FromName:    "John Doe",
		FromEmail:   "john.doe@example.com",
		Body:        "This is a preview of your Laravel-style email template. The styling matches Laravel's default mail templates exactly!",
	}

	htmlBody, err := cs.MailRendererService.RenderContactEmail(emailData)
	if err != nil {
		return "", err
	}
	return htmlBody, nil
}

// You can move this helper here as well
func generateFallbackHTML(req ContactRequest) string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Contact Form Message</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2e7d32;">New Contact Form Message</h2>
        <p><strong>From:</strong> ` + req.From + `</p>
        <p><strong>Subject:</strong> ` + req.Subject + `</p>
        <div style="background: #f6f6f6; padding: 15px; border-left: 4px solid #2e7d32; margin: 20px 0;">
            <p><strong>Message:</strong></p>
            <p>` + strings.ReplaceAll(req.Body, "\n", "<br>") + `</p>
        </div>
    </div>
</body>
</html>`
}
