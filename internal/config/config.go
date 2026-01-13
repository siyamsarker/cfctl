package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/siyamsarker/cfctl/pkg/cloudflare"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Version  int                  `yaml:"version" mapstructure:"version"`
	Defaults DefaultSettings      `yaml:"defaults" mapstructure:"defaults"`
	API      APISettings          `yaml:"api" mapstructure:"api"`
	UI       UISettings           `yaml:"ui" mapstructure:"ui"`
	Cache    CacheSettings        `yaml:"cache" mapstructure:"cache"`
	Accounts []cloudflare.Account `yaml:"accounts" mapstructure:"accounts"`
}

// DefaultSettings holds default application settings
type DefaultSettings struct {
	Account string `yaml:"account" mapstructure:"account"`
	Theme   string `yaml:"theme" mapstructure:"theme"`
	Output  string `yaml:"output" mapstructure:"output"`
}

// APISettings holds API configuration
type APISettings struct {
	Timeout int `yaml:"timeout" mapstructure:"timeout"`
	Retries int `yaml:"retries" mapstructure:"retries"`
}

// UISettings holds UI configuration
type UISettings struct {
	Confirmations bool `yaml:"confirmations" mapstructure:"confirmations"`
	Animations    bool `yaml:"animations" mapstructure:"animations"`
	Colors        bool `yaml:"colors" mapstructure:"colors"`
}

// CacheSettings holds cache configuration
type CacheSettings struct {
	DomainsTTL int  `yaml:"domains_ttl" mapstructure:"domains_ttl"`
	Enabled    bool `yaml:"enabled" mapstructure:"enabled"`
}

// Load loads configuration from file
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("get config path: %w", err)
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Set defaults
	setDefaults()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		return createDefaultConfig(configPath)
	}

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("get config path: %w", err)
	}

	viper.Set("version", c.Version)
	viper.Set("defaults", c.Defaults)
	viper.Set("api", c.API)
	viper.Set("ui", c.UI)
	viper.Set("cache", c.Cache)
	viper.Set("accounts", c.Accounts)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// AddAccount adds a new account to the configuration
func (c *Config) AddAccount(account cloudflare.Account) error {
	// Check if account with same name exists
	for i, acc := range c.Accounts {
		if acc.Name == account.Name {
			// Update existing account
			account.UpdatedAt = time.Now()
			c.Accounts[i] = account
			return c.Save()
		}
	}

	// Add new account
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	c.Accounts = append(c.Accounts, account)

	// Set as default if it's the first account
	if len(c.Accounts) == 1 {
		c.Accounts[0].Default = true
		c.Defaults.Account = account.Name
	}

	return c.Save()
}

// RemoveAccount removes an account from the configuration
func (c *Config) RemoveAccount(name string) error {
	for i, acc := range c.Accounts {
		if acc.Name == name {
			c.Accounts = append(c.Accounts[:i], c.Accounts[i+1:]...)

			// If this was the default account, set a new default
			if acc.Default && len(c.Accounts) > 0 {
				c.Accounts[0].Default = true
				c.Defaults.Account = c.Accounts[0].Name
			}

			return c.Save()
		}
	}
	return fmt.Errorf("account not found: %s", name)
}

// GetAccount retrieves an account by name
func (c *Config) GetAccount(name string) (*cloudflare.Account, error) {
	for _, acc := range c.Accounts {
		if acc.Name == name {
			return &acc, nil
		}
	}
	return nil, fmt.Errorf("account not found: %s", name)
}

// GetDefaultAccount returns the default account
func (c *Config) GetDefaultAccount() (*cloudflare.Account, error) {
	for _, acc := range c.Accounts {
		if acc.Default {
			return &acc, nil
		}
	}

	if len(c.Accounts) > 0 {
		return &c.Accounts[0], nil
	}

	return nil, fmt.Errorf("no accounts configured")
}

// SetDefaultAccount sets an account as the default
func (c *Config) SetDefaultAccount(name string) error {
	found := false
	for i := range c.Accounts {
		if c.Accounts[i].Name == name {
			c.Accounts[i].Default = true
			c.Defaults.Account = name
			found = true
		} else {
			c.Accounts[i].Default = false
		}
	}

	if !found {
		return fmt.Errorf("account not found: %s", name)
	}

	return c.Save()
}

func getConfigPath() (string, error) {
	// Check environment variable
	if path := os.Getenv("CFCTL_CONFIG"); path != "" {
		return path, nil
	}

	// Use ~/.config/cfctl/config.yaml
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "cfctl", "config.yaml"), nil
}

func setDefaults() {
	viper.SetDefault("version", 1)
	viper.SetDefault("defaults.theme", "dark")
	viper.SetDefault("defaults.output", "interactive")
	viper.SetDefault("api.timeout", 30)
	viper.SetDefault("api.retries", 3)
	viper.SetDefault("ui.confirmations", true)
	viper.SetDefault("ui.animations", true)
	viper.SetDefault("ui.colors", true)
	viper.SetDefault("cache.domains_ttl", 300)
	viper.SetDefault("cache.enabled", true)
}

func createDefaultConfig(path string) (*Config, error) {
	// Create directory if not exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	// Write default config
	if err := viper.SafeWriteConfigAs(path); err != nil {
		return nil, fmt.Errorf("write default config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal default config: %w", err)
	}

	return &cfg, nil
}
