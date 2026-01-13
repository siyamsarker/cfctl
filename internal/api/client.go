package api

import (
	"context"
	"fmt"
	"time"

	cf "github.com/cloudflare/cloudflare-go"
)

// Client wraps the Cloudflare API client
type Client struct {
	api     *cf.API
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
	var api *cf.API
	var err error

	if cfg.APIToken != "" {
		api, err = cf.NewWithAPIToken(cfg.APIToken)
	} else if cfg.APIKey != "" && cfg.Email != "" {
		api, err = cf.New(cfg.APIKey, cfg.Email)
	} else {
		return nil, fmt.Errorf("either API token or API key with email must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("create cloudflare client: %w", err)
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
	_, err := c.api.VerifyAPIToken(ctx)
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
