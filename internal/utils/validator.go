package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateURL validates a URL format
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL is required")
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	if u.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

// ValidateURLs validates multiple URLs
func ValidateURLs(urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("at least one URL is required")
	}

	if len(urls) > 30 {
		return fmt.Errorf("maximum 30 URLs allowed per request")
	}

	for i, urlStr := range urls {
		if err := ValidateURL(urlStr); err != nil {
			return fmt.Errorf("URL %d: %w", i+1, err)
		}
	}

	return nil
}

// ValidateHostname validates a hostname format
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.TrimPrefix(hostname, "https://")

	if strings.Contains(hostname, "/") {
		return fmt.Errorf("hostname should not contain path")
	}

	if strings.Contains(hostname, "?") {
		return fmt.Errorf("hostname should not contain query parameters")
	}

	if len(hostname) == 0 || len(hostname) > 253 {
		return fmt.Errorf("invalid hostname length")
	}

	return nil
}

// ValidateHostnames validates multiple hostnames
func ValidateHostnames(hostnames []string) error {
	if len(hostnames) == 0 {
		return fmt.Errorf("at least one hostname is required")
	}

	for i, hostname := range hostnames {
		if err := ValidateHostname(hostname); err != nil {
			return fmt.Errorf("hostname %d: %w", i+1, err)
		}
	}

	return nil
}

// ValidatePrefix validates a URL prefix
func ValidatePrefix(prefix string) error {
	if prefix == "" {
		return fmt.Errorf("prefix is required")
	}

	if err := ValidateURL(prefix); err != nil {
		return fmt.Errorf("prefix must be a valid URL: %w", err)
	}

	return nil
}

// ValidatePrefixes validates multiple prefixes
func ValidatePrefixes(prefixes []string) error {
	if len(prefixes) == 0 {
		return fmt.Errorf("at least one prefix is required")
	}

	for i, prefix := range prefixes {
		if err := ValidatePrefix(prefix); err != nil {
			return fmt.Errorf("prefix %d: %w", i+1, err)
		}
	}

	return nil
}

// ValidateTags validates cache tags
func ValidateTags(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	if len(tags) > 30 {
		return fmt.Errorf("maximum 30 tags allowed per request")
	}

	for i, tag := range tags {
		if tag == "" {
			return fmt.Errorf("tag %d: tag cannot be empty", i+1)
		}
	}

	return nil
}

// ParseCommaSeparated parses a comma-separated string into a slice
func ParseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
