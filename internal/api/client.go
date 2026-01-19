package api

import (
	"context"
	"fmt"
	"time"

	cfv6 "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

// Client wraps the Cloudflare API client
type Client struct {
	api     *cfv6.Client
	timeout time.Duration
	retries int
}

// ClientConfig holds configuration for creating a client
type ClientConfig struct {
	APIToken string
	APIKey   string
	Email    string
	Timeout  int
	Retries  int
}

// NewClient creates a new Cloudflare API client
func NewClient(cfg ClientConfig) (*Client, error) {
	var api *cfv6.Client

	if cfg.APIToken != "" {
		api = cfv6.NewClient(
			option.WithAPIToken(cfg.APIToken),
		)
	} else if cfg.APIKey != "" && cfg.Email != "" {
		api = cfv6.NewClient(
			option.WithAPIKey(cfg.APIKey),
			option.WithAPIEmail(cfg.Email),
		)
	} else {
		return nil, fmt.Errorf("either API token or API key with email must be provided")
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	if cfg.Timeout == 0 {
		timeout = 30 * time.Second
	}

	retries := cfg.Retries
	if retries == 0 {
		retries = 3
	}

	return &Client{
		api:     api,
		timeout: timeout,
		retries: retries,
	}, nil
}

// VerifyToken verifies the API credentials
func (c *Client) VerifyToken(ctx context.Context) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Try to verify the token by making a simple API call
	_, err := c.api.User.Get(ctx)
	if err != nil {
		errMsg := err.Error()

		// Provide helpful error messages for common issues
		if contains(errMsg, "code\":9109") || contains(errMsg, "Cannot use the access token from location") {
			return fmt.Errorf("IP restriction error: Your API token has IP address restrictions configured in Cloudflare. Please either remove the IP restrictions from your token or add your current IP address to the allowed list")
		}
		if contains(errMsg, "403") || contains(errMsg, "Forbidden") {
			return fmt.Errorf("authentication failed: Invalid API token or insufficient permissions. Please check your token has the required permissions")
		}
		if contains(errMsg, "401") || contains(errMsg, "Unauthorized") {
			return fmt.Errorf("authentication failed: Invalid API credentials. Please verify your token is correct")
		}

		return fmt.Errorf("verify credentials: %w", err)
	}

	return nil
}

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetTimeout returns the configured timeout
func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}

// GetRetries returns the configured retry count
func (c *Client) GetRetries() int {
	return c.retries
}
