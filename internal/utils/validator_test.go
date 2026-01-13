package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid http URL",
			url:     "http://example.com/page",
			wantErr: false,
		},
		{
			name:    "valid https URL",
			url:     "https://example.com/page",
			wantErr: false,
		},
		{
			name:    "invalid URL - no scheme",
			url:     "example.com/page",
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong scheme",
			url:     "ftp://example.com/page",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		wantErr  bool
	}{
		{
			name:     "valid hostname",
			hostname: "example.com",
			wantErr:  false,
		},
		{
			name:     "valid subdomain",
			hostname: "www.example.com",
			wantErr:  false,
		},
		{
			name:     "valid hostname - strips scheme",
			hostname: "http://example.com",
			wantErr:  false,
		},
		{
			name:     "invalid hostname - with path",
			hostname: "example.com/path",
			wantErr:  true,
		},
		{
			name:     "empty hostname",
			hostname: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHostname(tt.hostname)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePrefix(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		wantErr bool
	}{
		{
			name:    "valid prefix with scheme and path",
			prefix:  "http://example.com/blog",
			wantErr: false,
		},
		{
			name:    "valid prefix with https",
			prefix:  "https://www.example.com/api",
			wantErr: false,
		},
		{
			name:    "invalid prefix - no scheme",
			prefix:  "example.com/blog",
			wantErr: true,
		},
		{
			name:    "valid prefix without path",
			prefix:  "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty prefix",
			prefix:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePrefix(tt.prefix)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParsCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single item",
			input:    "example.com",
			expected: []string{"example.com"},
		},
		{
			name:     "multiple items",
			input:    "example.com, test.com, demo.com",
			expected: []string{"example.com", "test.com", "demo.com"},
		},
		{
			name:     "with extra whitespace",
			input:    "  example.com  ,  test.com  ",
			expected: []string{"example.com", "test.com"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only commas and spaces",
			input:    " , , ",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCommaSeparated(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
