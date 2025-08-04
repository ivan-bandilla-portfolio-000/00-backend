# Portfolio Backend with Laravel-Style Email Templates

This Go backend provides email functionality with Laravel-style email templates that look identical to Laravel's default mail styling.

## 🚀 Quick Start

1. **Copy environment variables:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your settings:**
   - Configure SMTP settings for email sending
   - Set your portfolio contact email
   - Update app name and URL

3. **Run the server:**
   ```bash
   go run main.go
   ```

## 📧 Email Templates

The backend uses Laravel-style email templates located in `templates/email/components/ui/`:

- **Layout Template**: `layout.tmpl` - Main email structure
- **Header Template**: `header.tmpl` - Email header with logo/branding  
- **Footer Template**: `footer.tmpl` - Email footer with copyright
- **Message Template**: `message.tmpl` - Contact form message content
- **Components**: `button.tmpl`, `panel.tmpl`, `table.tmpl`, `subcopy.tmpl`

## 🔗 API Endpoints

### Send Email
- **POST** `/send-email`
- **Body:**
  ```json
  {
    "from": "john@example.com",
    "name": "John Doe",
    "subject": "Contact Form Message", 
    "body": "Your message content here"
  }
  ```

### Preview Email Template
- **GET** `/preview-email`
- Shows how the Laravel-style email template looks
- Perfect for testing and design verification

### Discord Webhook (Optional)
- **POST** `/discord-webhook`
- Send notifications to Discord

## ✨ Features

- ✅ **Laravel-identical styling** - Emails look exactly like Laravel's default templates
- ✅ **Responsive design** - Works on all devices and email clients
- ✅ **Markdown support** - Content is parsed from Markdown to HTML
- ✅ **Component-based** - Modular template system
- ✅ **Fallback HTML** - Graceful degradation if templates fail
- ✅ **Environment configuration** - Easy setup via .env file
- ✅ **CORS enabled** - Ready for frontend integration

## 🛠 Development

### Preview Email Templates
Visit `http://localhost:8080/preview-email` to see the Laravel-style email template in action.

### Template Structure
```
templates/email/
├── components/ui/
│   ├── layout.tmpl      # Main layout
│   ├── header.tmpl      # Header component
│   ├── footer.tmpl      # Footer component
│   ├── message.tmpl     # Message content
│   ├── button.tmpl      # CTA buttons
│   ├── panel.tmpl       # Content panels
│   └── table.tmpl       # Data tables
└── themes/
    └── default.css      # Laravel's CSS
```

### Customization
- Edit templates in `templates/email/components/ui/`
- Modify styling in `templates/email/themes/default.css`
- Update email data structure in `services/email/service.go`

## 📝 Example Usage

```go
// Initialize email service
emailService, _ := email.NewEmailService()

// Create email data
data := email.EmailData{
    AppName:     "Your Portfolio",
    HeaderTitle: "Your Name",
    FromName:    "Contact Person",
    FromEmail:   "contact@example.com", 
    Body:        "Message content",
}

// Render Laravel-style email
html, _ := emailService.RenderContactEmail(data)
```

The rendered email will be styled exactly like Laravel's default mail templates with professional appearance and cross-client compatibility.
