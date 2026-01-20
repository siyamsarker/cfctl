package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "Valid Token",
			config: ClientConfig{
				APIToken: "test-token",
			},
			wantErr: false,
		},
		{
			name: "Valid Key and Email",
			config: ClientConfig{
				APIKey: "test-key",
				Email:  "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Missing Credentials",
			config: ClientConfig{
				APIToken: "",
				APIKey:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestClientTimeout(t *testing.T) {
	cfg := ClientConfig{
		APIToken: "test-token",
		Timeout:  10,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Second, client.timeout)
}

func TestClientDefaultTimeout(t *testing.T) {
	cfg := ClientConfig{
		APIToken: "test-token",
		Timeout:  0, // Should default to 30
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.Equal(t, 30*time.Second, client.timeout)
}

func TestClientRetries(t *testing.T) {
	cfg := ClientConfig{
		APIToken: "test-token",
		Retries:  5,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.Equal(t, 5, client.retries)
}

func TestClientDefaultRetries(t *testing.T) {
	cfg := ClientConfig{
		APIToken: "test-token",
		Retries:  0, // Should default to 3
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.Equal(t, 3, client.retries)
}
