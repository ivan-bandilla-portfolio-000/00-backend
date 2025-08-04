# Laravel-Style Email Templates for Go

This email template system replicates Laravel's mail template structure for Go applications.

## Directory Structure

```
templates/email/
├── components/
│   └── ui/
│       ├── button.tmpl      # Button component
│       ├── footer.tmpl      # Footer component  
│       ├── header.tmpl      # Header component
│       ├── layout.tmpl      # Main layout template
│       ├── message.tmpl     # Message content template
│       ├── panel.tmpl       # Panel component
│       ├── subcopy.tmpl     # Subcopy component
│       └── table.tmpl       # Table component
├── themes/
│   └── default.css          # Laravel's default email styles
├── body-notification.tmpl   # Notification email template
└── markdown.go              # Markdown parser (Laravel-style)
```

## Usage

### Basic Setup

```go
package main

import (
    "log"
    "./services/email"
)

func main() {
    emailService, err := email.NewEmailService()
    if err != nil {
        log.Fatal("Failed to initialize email service:", err)
    }
    
    data := email.EmailData{
        AppName:     "Your App Name",
        Subject:     "Email Subject",
        SiteURL:     "https://yoursite.com",
        HeaderTitle: "Your Brand",
        FromName:    "John Doe",
        FromEmail:   "john@example.com",
        Body:        "Your email content here",
    }
    
    html, err := emailService.RenderContactEmail(data)
    if err != nil {
        log.Fatal("Failed to render email:", err)
    }
    
    // Use the rendered HTML for sending emails
    log.Println(html)
}
```

### Template Components

#### Layout (`layout.tmpl`)
Main template that wraps all email content with proper HTML structure and styling.

#### Header (`header.tmpl`)
Renders the email header with logo/brand name and URL.

#### Footer (`footer.tmpl`)
Renders the email footer with copyright and company information.

#### Button (`button.tmpl`)
Creates styled buttons for CTAs:
- Supports `url`, `color` (primary, success, error), and `align` properties
- Follows Laravel's button styling

#### Panel (`panel.tmpl`)
Creates bordered content panels for highlighting information.

#### Table (`table.tmpl`)
Renders properly styled tables for data display.

#### Message (`message.tmpl`)
Template for contact form messages and notifications.

### Markdown Support

The system includes a Markdown parser that supports:
- Headers (`#`, `##`, `###`)
- Bold text (`**text**`)
- Line breaks and paragraphs
- Horizontal rules (`---`)

### Styling

Uses Laravel's exact default email CSS with:
- Responsive design
- Cross-client compatibility
- Proper color scheme
- Professional typography

### Features

- ✅ Laravel-compatible template structure
- ✅ Markdown parsing (similar to `Illuminate\Mail\Markdown::parse()`)
- ✅ Component-based architecture
- ✅ Responsive email design
- ✅ Cross-client compatibility
- ✅ Professional styling matching Laravel's defaults

### Example Output

The templates generate HTML emails that look identical to Laravel's default mail templates, with:
- Clean, professional appearance
- Proper mobile responsiveness
- Consistent branding
- Accessible markup

Run `go run example-email.go` to see a working example.
