package config

import (
	"fmt"
)

type VaultConfig struct {
	VaultAddress string `hcl:"address"`
	VaultToken   string `hcl:"token"`
	SSHConfig    SSH    `hcl:"ssh,block"`
}

type SSH struct {
	Username string `hcl:"username"`
	Password string `hcl:"password"`
	Port     uint   `hcl:"port"`
}

func (c *VaultConfig) Validate() error {
	if c.VaultAddress == "" {
		return fmt.Errorf("error with Vault address: %w", ErrBlankParam)
	}
	if c.VaultToken == "" {
		return fmt.Errorf("error with Vault token: %w", ErrBlankParam)
	}
	if c.SSHConfig.Username == "" {
		return fmt.Errorf("error with Vault SSH user: %w", ErrBlankParam)
	}
	if c.SSHConfig.Password == "" {
		return fmt.Errorf("error with Vault SSH password: %w", ErrBlankParam)
	}
	if c.SSHConfig.Port == 0 {
		return fmt.Errorf("error with Vault address: %w", ErrBlankParam)
	}
	return nil
}
