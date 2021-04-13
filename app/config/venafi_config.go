package config

import "fmt"

type VenafiCloudConnection struct {
	APIKey string `hcl:"apikey,optional"`
	Zone   string `hcl:"zone,optional"`
}

type VenafiTPPConnection struct {
	URL      string `hcl:"url"`
	Username string `hcl:"username"`
	// TODO: support access token
	Password string `hcl:"password"`
	Policy   string `hcl:"policy"`
}

func (c *VenafiCloudConnection) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("error with Venafi API key: %w", ErrBlankParam)
	}
	if c.Zone == "" {
		return fmt.Errorf("error with Venafi Zone: %w", ErrBlankParam)
	}
	return nil
}

func (c *VenafiTPPConnection) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("error with TPP URL: %w", ErrBlankParam)
	}
	if c.Policy == "" {
		return fmt.Errorf("error with TPP Policy: %w", ErrBlankParam)
	}
	if c.Username == "" {
		return fmt.Errorf("error with TPP Username: %w", ErrBlankParam)
	}
	if c.Password == "" {
		return fmt.Errorf("error with TPP Password: %w", ErrBlankParam)
	}
	return nil
}
