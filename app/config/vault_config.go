package config

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
)

type VaultConfig struct {
	VaultAddress string `hcl:"api_address"`
	VaultToken   string `hcl:"token"`
	SSHConfig    []SSH  `hcl:"ssh,block"`
}

type SSH struct {
	Hostname string `hcl:"hostname"`
	Username string `hcl:"username"`
	Password string `hcl:"password"`
	Port     uint   `hcl:"port"`
}

func (c *VaultConfig) Validate() error {
	if c.VaultAddress == "" {
		return fmt.Errorf("error with Vault address: %w", errors.ErrBlankParam)
	}
	if c.VaultToken == "" {
		return fmt.Errorf("error with Vault token: %w", errors.ErrBlankParam)
	}
	for _, ssh := range c.SSHConfig {
		if ssh.Hostname == "" {
			return fmt.Errorf("error with Vault SSH Hostname: %w", errors.ErrBlankParam)
		}
		if ssh.Username == "" {
			return fmt.Errorf("error with Vault SSH user: %w", errors.ErrBlankParam)
		}
		if ssh.Password == "" {
			return fmt.Errorf("error with Vault SSH password: %w", errors.ErrBlankParam)
		}
		if ssh.Port == 0 {
			return fmt.Errorf("error with Vault address: %w", errors.ErrBlankParam)
		}
	}
	return nil
}
