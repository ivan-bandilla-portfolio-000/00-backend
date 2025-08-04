package email

import (
	"html/template"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

// MarkdownToHTMLUnsafe converts markdown to HTML or passes through HTML content WITHOUT sanitization
func MarkdownToHTMLUnsafe(content string) template.HTML {
	// If content is already HTML, return it as-is without sanitization
	if isHTML(content) {
		return template.HTML(content)
	}

	// Convert markdown to HTML using gomarkdown library
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	// Create HTML renderer with safe settings
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert markdown to HTML using the renderer
	renderedHTML := markdown.Render(doc, renderer)

	// Return unsanitized HTML
	return template.HTML(string(renderedHTML))
}

// Template function map for Go templates
func GetTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"MarkdownToHTML": MarkdownToHTML,
		"default": func(defaultVal, val interface{}) interface{} {
			if val == nil || val == "" {
				return defaultVal
			}
			return val
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
	}
}

// isHTML checks if content contains HTML tags
func isHTML(content string) bool {
	re := regexp.MustCompile(`</?[a-z][\s\S]*>`)
	return re.MatchString(content)
}

// MarkdownToHTML converts markdown to HTML or passes through sanitized HTML content
func MarkdownToHTML(content string) template.HTML {
	// If content is already HTML, sanitize it
	if isHTML(content) {
		policy := bluemonday.UGCPolicy()
		return template.HTML(policy.Sanitize(content))
	}

	// Convert markdown to HTML using gomarkdown library
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	// Create HTML renderer with safe settings
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert markdown to HTML using the renderer
	renderedHTML := markdown.Render(doc, renderer)

	// Sanitize the generated HTML
	policy := bluemonday.UGCPolicy()
	sanitizedHTML := policy.Sanitize(string(renderedHTML)) // Convert []byte to string here

	return template.HTML(sanitizedHTML)
}
