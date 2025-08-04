package email

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// EmailData represents the data structure for email templates
type EmailData struct {
	AppName     string
	Subject     string
	SiteURL     string
	LogoURL     string
	HeaderTitle string
	Year        int
	FromName    string
	FromEmail   string
	Body        string
	Slot        string
	Header      template.HTML
	Footer      template.HTML
	Subcopy     template.HTML
}

// ComponentData represents data for email components
type ComponentData struct {
	Url   string
	Slot  string
	Color string
	Align string
}

// EmailService handles email template rendering
type EmailService struct {
	templates *template.Template
}

// NewEmailService creates a new email service with loaded templates
func NewEmailService() (*EmailService, error) {
	// Use the function map from markdown.go
	funcMap := GetTemplateFuncMap()

	// Parse all email templates
	templateDir := "./templates/email"
	tmpl := template.New("").Funcs(funcMap)

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".tmpl") || strings.HasSuffix(path, ".css") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Get relative path from templateDir
			relPath, err := filepath.Rel(templateDir, path)
			if err != nil {
				return err
			}

			// Convert Windows path separators to forward slashes for template names
			templateName := strings.ReplaceAll(relPath, "\\", "/")

			_, err = tmpl.New(templateName).Parse(string(content))
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &EmailService{templates: tmpl}, nil
}

// RenderHeader renders the email header component
func (es *EmailService) RenderHeader(url, slot string) (template.HTML, error) {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"url":  url,
		"slot": slot,
	}

	err := es.templates.ExecuteTemplate(&buf, "components/ui/header.tmpl", data)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// RenderFooter renders the email footer component
func (es *EmailService) RenderFooter(content string) (template.HTML, error) {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"slot": content,
	}

	err := es.templates.ExecuteTemplate(&buf, "components/ui/footer.tmpl", data)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}

// RenderContactEmail renders the contact form email
func (es *EmailService) RenderContactEmail(data EmailData) (string, error) {
	var buf bytes.Buffer

	// Set default values if not provided
	if data.Year == 0 {
		data.Year = time.Now().Year()
	}
	if data.HeaderTitle == "" {
		data.HeaderTitle = "Portfolio"
	}
	if data.AppName == "" {
		data.AppName = data.HeaderTitle
	}
	if data.SiteURL == "" {
		data.SiteURL = "https://yourdomain.com"
	}

	// Render header
	header, err := es.RenderHeader(data.SiteURL, data.HeaderTitle)
	if err != nil {
		return "", err
	}
	data.Header = header

	// Render footer
	footerContent := "© " + strconv.Itoa(data.Year) + " " + data.HeaderTitle + ". All rights reserved."
	footer, err := es.RenderFooter(footerContent)
	if err != nil {
		return "", err
	}
	data.Footer = footer

	// Render message content
	var msgBuf bytes.Buffer
	err = es.templates.ExecuteTemplate(&msgBuf, "components/ui/message.tmpl", data)
	if err != nil {
		return "", err
	}
	data.Slot = msgBuf.String()

	// Render layout
	err = es.templates.ExecuteTemplate(&buf, "components/ui/layout.tmpl", data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
