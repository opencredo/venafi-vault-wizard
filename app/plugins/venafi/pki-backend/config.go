package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
)

type VenafiPKIBackendConfig struct {
	// MountPath is not decoded directly by using the struct tags, and is instead populated by plugins.LookupPlugin
	// when it is initialised
	MountPath string
	// Version is not decoded directly by using the struct tags, and is instead populated by plugins.LookupPlugin
	// when it is initialised
	Version string

	Roles []Role `hcl:"role,block"`
}

type Role struct {
	Name      string               `hcl:"role,label"`
	Secret    venafi.VenafiSecret  `hcl:"secret,block"`
	TestCerts []CertificateRequest `hcl:"test_certificate,block"`
}

type CertificateRequest struct {
	CommonName string `hcl:"common_name"`
}

func (c *VenafiPKIBackendConfig) ValidateConfig() error {
	if len(c.Roles) == 0 {
		return fmt.Errorf("error at least one role must be provided: %w", errors.ErrBlankParam)
	}
	for _, role := range c.Roles {
		err := role.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Role) Validate() error {
	return r.Secret.Validate()
}
