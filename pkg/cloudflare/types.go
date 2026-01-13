package cloudflare

import "time"

// Zone represents a Cloudflare zone/domain
type Zone struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Plan   Plan   `json:"plan"`
}

// Plan represents a Cloudflare plan
type Plan struct {
	Name string `json:"name"`
}

// PurgeRequest represents cache purge request
type PurgeRequest struct {
	Files           []string `json:"files,omitempty"`
	Hosts           []string `json:"hosts,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Prefixes        []string `json:"prefixes,omitempty"`
	PurgeEverything bool     `json:"purge_everything,omitempty"`
}

// APIResponse represents standard Cloudflare API response
type APIResponse struct {
	Success  bool        `json:"success"`
	Errors   []APIError  `json:"errors"`
	Messages []string    `json:"messages"`
	Result   interface{} `json:"result"`
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Account represents stored account configuration
type Account struct {
	Name      string    `yaml:"name" mapstructure:"name"`
	Email     string    `yaml:"email" mapstructure:"email"`
	AuthType  string    `yaml:"auth_type" mapstructure:"auth_type"`
	Default   bool      `yaml:"default" mapstructure:"default"`
	CreatedAt time.Time `yaml:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at" mapstructure:"updated_at"`
}
