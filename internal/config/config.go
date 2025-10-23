package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	// ACME server settings
	Email   string `yaml:"email"`
	Server  string `yaml:"server"`
	Staging bool   `yaml:"staging"`

	// Storage settings
	CertDir string `yaml:"cert_dir"`
	
	// Default settings
	KeyType     string `yaml:"key_type"`
	RenewDays   int    `yaml:"renew_days"`
	
	// DNS provider settings
	DNSProviders map[string]map[string]string `yaml:"dns_providers"`
	
	// Deployment hooks
	DeployHooks map[string]DeployHook `yaml:"deploy_hooks"`
	
	// Notification settings
	Notifications NotificationConfig `yaml:"notifications"`
}

// DeployHook represents a deployment hook configuration
type DeployHook struct {
	Type     string            `yaml:"type"`
	Script   string            `yaml:"script"`
	Settings map[string]string `yaml:"settings"`
}

// NotificationConfig represents notification settings
type NotificationConfig struct {
	Enabled bool              `yaml:"enabled"`
	Type    string            `yaml:"type"` // email, webhook, slack, etc.
	Settings map[string]string `yaml:"settings"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	
	return &Config{
		Server:    "https://acme-v02.api.letsencrypt.org/directory",
		Staging:   false,
		CertDir:   filepath.Join(homeDir, ".acme-go"),
		KeyType:   "ec-256",
		RenewDays: 30,
		DNSProviders: make(map[string]map[string]string),
		DeployHooks: make(map[string]DeployHook),
		Notifications: NotificationConfig{
			Enabled: false,
		},
	}
}

// Load loads configuration from file or returns default config
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".acme-go.yaml")
	
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure cert directory exists
	if err := os.MkdirAll(config.CertDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cert directory: %w", err)
	}

	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".acme-go.yaml")
	
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDNSProvider returns DNS provider configuration
func (c *Config) GetDNSProvider(name string) (map[string]string, bool) {
	provider, exists := c.DNSProviders[name]
	return provider, exists
}

// SetDNSProvider sets DNS provider configuration
func (c *Config) SetDNSProvider(name string, settings map[string]string) {
	if c.DNSProviders == nil {
		c.DNSProviders = make(map[string]map[string]string)
	}
	c.DNSProviders[name] = settings
}

// GetDeployHook returns deployment hook configuration
func (c *Config) GetDeployHook(name string) (DeployHook, bool) {
	hook, exists := c.DeployHooks[name]
	return hook, exists
}

// SetDeployHook sets deployment hook configuration
func (c *Config) SetDeployHook(name string, hook DeployHook) {
	if c.DeployHooks == nil {
		c.DeployHooks = make(map[string]DeployHook)
	}
	c.DeployHooks[name] = hook
}