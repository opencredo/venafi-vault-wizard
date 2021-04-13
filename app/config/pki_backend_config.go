package config

import "fmt"

type PKIBackendConfig struct {
	VenafiSecret              string
	RoleName                  string
	TestCertificateCommonName string
}

func (c *PKIBackendConfig) Validate() error {
	if c.VenafiSecret == "" {
		return fmt.Errorf("error with Venafi secret: %w", ErrBlankParam)
	}
	if c.RoleName == "" {
		return fmt.Errorf("error with Venafi role name: %w", ErrBlankParam)
	}
	return nil
}
