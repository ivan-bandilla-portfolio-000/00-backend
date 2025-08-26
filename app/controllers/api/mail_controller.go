package api_controllers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"portfolio-backend/services/email"
	validation_email "portfolio-backend/services/validation/email"
	"portfolio-backend/services/recaptcha"
)

type EmailController struct {
	MailService *email.MailService
}

func NewEmailController(mailService *email.MailService) *EmailController {
	return &EmailController{MailService: mailService}
}

type EmailRequest struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Name    string `json:"name,omitempty"`
	RecaptchaToken string `json:"recaptchaToken,omitempty"`
}

// Handler: POST /send-email
func (ec *EmailController) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Determine client IP (prefer X-Forwarded-For / X-Real-IP)
    clientIP := func(r *http.Request) string {
        if ip := r.Header.Get("X-Real-IP"); ip != "" {
            return ip
        }
        if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
            parts := strings.Split(fwd, ",")
            return strings.TrimSpace(parts[0])
        }
        host, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            return r.RemoteAddr
        }
        return host
    }(r)

	clientUserAgent := r.Header.Get("User-Agent")

	// Verify reCAPTCHA v3 token (action must match client-side action)
    ok, score, err := recaptcha.Verify(req.RecaptchaToken, "contact_submit", 0.5, clientIP)
    if err != nil || !ok {
        log.Printf("recaptcha verification failed: ok=%v score=%v err=%v ip=%s user_agent=%q remote_addr=%s", ok, score, err, clientIP, clientUserAgent, r.RemoteAddr)
        http.Error(w, "reCAPTCHA verification failed", http.StatusForbidden)
        return
    }

	log.Printf("recaptcha verified: score=%.2f ip=%s user_agent=%q remote_addr=%s", score, clientIP, clientUserAgent, r.RemoteAddr)

	valid, err := validation_email.ValidateEmail(req.From)
	if err != nil {
		if vErr, ok := err.(*validation_email.ValidationError); ok && vErr.Code == "invalid_format" {
			http.Error(w, vErr.Message, http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to validate email", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Email is invalid, disposable, or does not exist", http.StatusBadRequest)
		return
	}

	contactReq := email.ContactRequest{
		From:    req.From,
		Subject: req.Subject,
		Body:    req.Body,
		Name:    req.Name,
	}

	err = ec.MailService.SendContactEmail(contactReq)
	if err != nil {
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
	if ec.MailService == nil {
		log.Println("MailService is nil in PreviewEmail handler")
		http.Error(w, "Email service not initialized", http.StatusInternalServerError)
		return
	}

	htmlBody, err := ec.MailService.GeneratePreviewEmail()
	if err != nil {
		http.Error(w, "Failed to render email template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlBody))
}
