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

// MailRendererService handles email template rendering
type MailRendererService struct {
	templates *template.Template
}

// NewMailRendererService creates a new email service with loaded templates
func NewMailRendererService() (*MailRendererService, error) {
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

	return &MailRendererService{templates: tmpl}, nil
}

// RenderHeader renders the email header component
func (es *MailRendererService) RenderHeader(url, slot string) (template.HTML, error) {
	header, err := es.renderToString("components/ui/header.tmpl", map[string]interface{}{
		"url":  url,
		"slot": slot,
	})
	return template.HTML(header), err
}

// RenderFooter renders the email footer component
func (es *MailRendererService) RenderFooter(content string) (template.HTML, error) {
	footer, err := es.renderToString("components/ui/footer.tmpl", map[string]interface{}{
		"slot": content,
	})
	return template.HTML(footer), err
}

// RenderContactEmail renders the contact form email
func (es *MailRendererService) RenderContactEmail(data EmailData) (string, error) {
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

	footerContent := "© " + strconv.Itoa(data.Year) + " " + data.HeaderTitle + ". All rights reserved."
	footer, err := es.RenderFooter(footerContent)
	if err != nil {
		return "", err
	}
	data.Footer = footer

	// Render message content
	slot, err := es.renderToString("components/ui/message.tmpl", data)
	if err != nil {
		return "", err
	}
	data.Slot = slot

	// Render layout
	return es.renderToString("components/ui/layout.tmpl", data)
}

func (es *MailRendererService) renderToString(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := es.templates.ExecuteTemplate(&buf, name, data)
	return buf.String(), err
}

func (es *MailRendererService) Templates() *template.Template {
	return es.templates
}
