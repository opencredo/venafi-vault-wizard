package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/zclconf/go-cty/cty"
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

// WriteHCL uses the hclwrite package to encode itself into HCL. It supports $ENVVARS for the string values, in that
// format. This allows users in a wizard to specify the string params in a shell-like syntax, which will then be
// serialised into the HCL syntax of env("ENVVARS")
func (c *VaultConfig) WriteHCL(hclBody *hclwrite.Body) {
	vaultConfigBlock := hclBody.AppendNewBlock("vault", nil)
	vaultConfigBody := vaultConfigBlock.Body()

	WriteStringAttributeToHCL("api_address", c.VaultAddress, vaultConfigBody)
	WriteStringAttributeToHCL("token", c.VaultToken, vaultConfigBody)

	for _, sshHost := range c.SSHConfig {
		vaultConfigBody.AppendNewline()

		sshHostBlock := vaultConfigBody.AppendNewBlock("ssh", nil)
		sshHostBody := sshHostBlock.Body()

		WriteStringAttributeToHCL("hostname", sshHost.Hostname, sshHostBody)
		WriteStringAttributeToHCL("username", sshHost.Username, sshHostBody)
		WriteStringAttributeToHCL("password", sshHost.Password, sshHostBody)
		sshHostBody.SetAttributeValue("port", cty.NumberUIntVal(uint64(sshHost.Port)))
	}
}
