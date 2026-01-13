package utils

import (
	"fmt"
	"strings"
)

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatList formats a slice of strings as a numbered list
func FormatList(items []string) string {
	if len(items) == 0 {
		return "(none)"
	}

	var sb strings.Builder
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}
	return sb.String()
}

// PluralizeWord returns singular or plural form based on count
func PluralizeWord(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// FormatCount formats a count with proper pluralization
func FormatCount(count int, singular, plural string) string {
	return fmt.Sprintf("%d %s", count, PluralizeWord(count, singular, plural))
}
