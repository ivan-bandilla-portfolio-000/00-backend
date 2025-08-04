package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"portfolio-backend/services/email"
	"portfolio-backend/utils"

	"gopkg.in/gomail.v2"
)

type EmailController struct {
	EmailService *email.EmailService
}

func NewEmailController(emailService *email.EmailService) *EmailController {
	return &EmailController{
		EmailService: emailService,
	}
}

type EmailRequest struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
}

// Handler: POST /send-email
func (ec *EmailController) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	emailData := email.EmailData{
		AppName:     utils.GetEnvOrDefault("APP_NAME", "Portfolio"),
		Subject:     req.Subject,
		SiteURL:     utils.GetEnvOrDefault("APP_URL", "https://yourdomain.com"),
		HeaderTitle: utils.GetEnvOrDefault("APP_NAME", "Portfolio"),
		Year:        time.Now().Year(),
		FromName:    utils.GetNameFromEmail(req.From, req.Name),
		FromEmail:   req.From,
		Body:        req.Body,
	}

	htmlBody, err := ec.EmailService.RenderContactEmail(emailData)
	if err != nil {
		log.Printf("Failed to render email template: %v", err)
		htmlBody = generateFallbackHTML(req)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", req.From)
	m.SetHeader("To", os.Getenv("THIS_PORTFOLIO_CONTACT_EMAIL"))
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", htmlBody)
	plainBody := utils.StripHTMLTags(htmlBody)
	m.AddAlternative("text/plain", plainBody)

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		http.Error(w, "Invalid SMTP_PORT", http.StatusInternalServerError)
		return
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your message was sent successfully!",
	})
}

// Handler: GET /preview-email
func (ec *EmailController) PreviewEmail(w http.ResponseWriter, r *http.Request) {
	if ec.EmailService == nil {
		http.Error(w, "Email service not initialized", http.StatusInternalServerError)
		return
	}

	emailData := email.EmailData{
		AppName:     utils.GetEnvOrDefault("APP_NAME", "Ivan Bandilla Portfolio"),
		Subject:     "Contact Form Preview",
		SiteURL:     utils.GetEnvOrDefault("APP_URL", "https://ivanbandilla.dev"),
		HeaderTitle: utils.GetEnvOrDefault("APP_NAME", "Ivan Bandilla"),
		Year:        time.Now().Year(),
		FromName:    "John Doe",
		FromEmail:   "john.doe@example.com",
		Body:        "This is a preview of your Laravel-style email template. The styling matches Laravel's default mail templates exactly!",
	}

	htmlBody, err := ec.EmailService.RenderContactEmail(emailData)
	if err != nil {
		http.Error(w, "Failed to render email template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlBody))
}

// --- Helpers ---

func generateFallbackHTML(req EmailRequest) string {
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
