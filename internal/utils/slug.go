package utils

import (
	"strings"
	"unicode"
)

// GenerateSlug converts a name into URL-friendly slug
// Example: "Woman Ready to Wear" → "woman-ready-to-wear"
func GenerateSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))

	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('-')
		}
		// skip other chars
	}

	out := b.String()

	// Collapse multiple hyphens
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}

	// Trim leading/trailing hyphens
	out = strings.Trim(out, "-")

	return out
}