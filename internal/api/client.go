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
		return fmt.Errorf("verify credentials: %w", err)
	}

	return nil
}

// GetTimeout returns the configured timeout
func (c *Client) GetTimeout() time.Duration {
	return c.timeout
}

// GetRetries returns the configured retry count
func (c *Client) GetRetries() int {
	return c.retries
}
