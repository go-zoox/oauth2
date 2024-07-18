package oauth2

import (
	"fmt"

	"github.com/go-zoox/core-utils/safe"
)

var clients = safe.NewMap[string, any]()

// Register registers a new oauth2 service provider.
func Register(provider string, cfg *Config) error {
	if provider == "" {
		return fmt.Errorf("oauth2: provider is empty")
	}

	if clients.Has(provider) {
		return fmt.Errorf("oauth2: provider(%s) already registered", provider)
	}

	if cfg.Name == "" {
		cfg.Name = provider
	}

	clients.Set(provider, cfg)
	return nil
}

// Get gets the oauth2 service provider by name.
func Get(provider string) (*Config, error) {
	if provider == "" {
		return nil, fmt.Errorf("oauth2: provider is empty")
	}

	if !clients.Has(provider) {
		return nil, fmt.Errorf("oauth2: provider(%s) not registered", provider)
	}

	return clients.Get(provider).(*Config), nil
}
