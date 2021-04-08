package config

import (
	"fmt"
)

type VaultConfig struct {
	VaultAddress string
	VaultToken   string
	SSHUser      string
	SSHPassword  string
	SSHPort      uint
	MountPath    string
}

func (c *VaultConfig) Validate() error {
	if c.VaultAddress == "" {
		return fmt.Errorf("error with Vault address: %w", ErrBlankParam)
	}
	if c.VaultToken == "" {
		return fmt.Errorf("error with Vault token: %w", ErrBlankParam)
	}
	if c.SSHUser == "" {
		return fmt.Errorf("error with Vault SSH user: %w", ErrBlankParam)
	}
	if c.SSHPassword == "" {
		return fmt.Errorf("error with Vault SSH password: %w", ErrBlankParam)
	}
	if c.SSHPort == 0 {
		return fmt.Errorf("error with Vault address: %w", ErrBlankParam)
	}
	if c.MountPath == "" {
		return fmt.Errorf("error with backend mount path: %w", ErrBlankParam)
	}
	return nil
}
