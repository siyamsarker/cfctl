package config

import (
	"fmt"
	"net/mail"
	"strings"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateAccountName validates an account name
func ValidateAccountName(name string) error {
	if name == "" {
		return fmt.Errorf("account name is required")
	}

	if len(name) < 3 {
		return fmt.Errorf("account name must be at least 3 characters")
	}

	if len(name) > 50 {
		return fmt.Errorf("account name must be less than 50 characters")
	}

	// Allow alphanumeric, dash, underscore, dot, and space
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == ' ') {
			return fmt.Errorf("account name can only contain letters, numbers, dashes, underscores, dots, and spaces")
		}
	}

	return nil
}

// ValidateAPIToken validates an API token format
func ValidateAPIToken(token string) error {
	if token == "" {
		return fmt.Errorf("API token is required")
	}

	token = strings.TrimSpace(token)

	if len(token) < 40 {
		return fmt.Errorf("API token appears to be too short")
	}

	return nil
}

// ValidateAPIKey validates a global API key format
func ValidateAPIKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key is required")
	}

	key = strings.TrimSpace(key)

	if len(key) < 32 {
		return fmt.Errorf("API key appears to be too short")
	}

	return nil
}
