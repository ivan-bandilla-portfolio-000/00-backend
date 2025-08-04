package utils

import (
	"regexp"
	"strings"
)

// Converts an email to a display name if name is empty
func GetNameFromEmail(email, name string) string {
	if name != "" {
		return name
	}
	if parts := strings.Split(email, "@"); len(parts) > 0 {
		return strings.Title(strings.ReplaceAll(parts[0], ".", " "))
	}
	return email
}

// Removes HTML tags from a string
func StripHTMLTags(html string) string {
	re := regexp.MustCompile(`<.*?>`)
	return strings.TrimSpace(re.ReplaceAllString(html, ""))
}
