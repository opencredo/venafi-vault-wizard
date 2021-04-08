package config

import "fmt"

type VenafiConnectionConfig interface {
	GetAsMap() map[string]interface{}
}

type VenafiCloudConfig struct {
	VenafiAPIKey string
	VenafiZone   string
}

func (c *VenafiCloudConfig) Validate() error {
	if c.VenafiAPIKey == "" {
		return fmt.Errorf("error with Venafi API key: %w", ErrBlankParam)
	}
	if c.VenafiZone == "" {
		return fmt.Errorf("error with Venafi Zone: %w", ErrBlankParam)
	}
	return nil
}

func (c *VenafiCloudConfig) GetAsMap() map[string]interface{} {
	return map[string]interface{}{
		"apikey": c.VenafiAPIKey,
		"zone":   c.VenafiZone,
	}
}

type VenafiTPPConfig struct {
	URL string
	VenafiZone string
	Username string
	// TODO: only allow access token, not password
	Password string
}

func (c *VenafiTPPConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("error with TPP URL: %w", ErrBlankParam)
	}
	if c.VenafiZone == "" {
		return fmt.Errorf("error with TPP URL: %w", ErrBlankParam)
	}
	if c.Username == "" {
		return fmt.Errorf("error with TPP URL: %w", ErrBlankParam)
	}
	if c.Password == "" {
		return fmt.Errorf("error with TPP URL: %w", ErrBlankParam)
	}
	return nil
}

func (c *VenafiTPPConfig) GetAsMap() map[string]interface{} {
	return map[string]interface{}{
		"url":          c.URL,
		"zone":         c.VenafiZone,
		"tpp_user":     c.Username,
		"tpp_password": c.Password,
	}
}
