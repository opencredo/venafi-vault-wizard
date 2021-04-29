package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
)

type VenafiPKIBackendConfig struct {
	// MountPath is not decoded directly by using the struct tags, and is instead populated by
	// ParseVenafiPKIBackendConfig when it is initialised
	MountPath string
	// Version is not decoded directly by using the struct tags, and is instead populated by ParseVenafiPKIBackendConfig
	// when it is initialised
	Version string

	Roles []Role `hcl:"role,block"`
}

type Role struct {
	Name      string                      `hcl:"role,label"`
	Zone      string                      `hcl:"zone,optional"`
	Secret    venafi.VenafiSecret         `hcl:"secret,block"`
	TestCerts []venafi.CertificateRequest `hcl:"test_certificate,block"`
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
	err := r.Secret.Validate()
	if err != nil {
		return err
	}

	// Zone is optional for pki-monitor but required for pki-backend so check it here
	if zone, ok := r.Secret.GetAsMap()["zone"]; !ok || zone == "" {
		return fmt.Errorf("error zone must be specified in venafi-pki-backend secret")
	}

	return nil
}
