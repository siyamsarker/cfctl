package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email - no @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email - no domain",
			email:   "user@",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAccountName(t *testing.T) {
	tests := []struct {
		name    string
		accName string
		wantErr bool
	}{
		{
			name:    "valid name",
			accName: "my-account",
			wantErr: false,
		},
		{
			name:    "valid name with underscores",
			accName: "my_account_123",
			wantErr: false,
		},
		{
			name:    "empty name",
			accName: "",
			wantErr: true,
		},
		{
			name:    "name too long",
			accName: "this-is-a-very-long-account-name-that-exceeds-fifty-chars",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			accName: "account@name",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountName(tt.accName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAPIToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token length",
			token:   "abcdef1234567890abcdef1234567890abcdef1234567890",
			wantErr: false,
		},
		{
			name:    "token too short",
			token:   "short",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIToken(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
