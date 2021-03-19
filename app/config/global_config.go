package config

import (
	"errors"
	"fmt"
)

var ErrBlankParam = errors.New("cannot be blank")

type GlobalConfig struct {
	VaultAddress string
	VaultToken   string
	SSHUser      string
	SSHPassword  string
	SSHPort      uint
	VenafiAPIKey string
	VenafiZone   string
	MountPath    string
}

func (c *GlobalConfig) Validate() error {
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
	if c.VenafiAPIKey == "" {
		return fmt.Errorf("error with Venafi API key: %w", ErrBlankParam)
	}
	if c.VenafiZone == "" {
		return fmt.Errorf("error with Venafi Zone: %w", ErrBlankParam)
	}
	if c.MountPath == "" {
		return fmt.Errorf("error with backend mount path: %w", ErrBlankParam)
	}
	return nil
}
