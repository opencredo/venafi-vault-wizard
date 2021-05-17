package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
)

type VenafiPKIMonitorConfig struct {
	// MountPath is not decoded directly by using the struct tags, and is instead populated by plugins.LookupPlugin
	// when it is initialised
	MountPath string
	// Version is not decoded directly by using the struct tags, and is instead populated by plugins.LookupPlugin
	// when it is initialised
	Version string

	Role Role `hcl:"role,block"`
}

type Role struct {
	Name string `hcl:"role,label"`

	Secret venafi.VenafiSecret `hcl:"secret,block"`

	EnforcementPolicy Policy  `hcl:"enforcement_policy,block"`
	ImportPolicy      *Policy `hcl:"import_policy,block"`

	IntermediateCert *venafi.CertificateRequest `hcl:"intermediate_certificate,block"`
	RootCert         *venafi.CertificateRequest `hcl:"root_certificate,block"`

	TestCerts []venafi.CertificateRequest `hcl:"test_certificate,block"`

	GenerateLease bool   `hcl:"generate_lease,optional"`
	AllowAnyName  bool   `hcl:"allow_any_name,optional"`
	TTL           string `hcl:"ttl,optional"`
	MaxTTL        string `hcl:"max_ttl,optional"`
}

type Policy struct {
	Zone string `hcl:"zone"`
}

func (c *VenafiPKIMonitorConfig) ValidateConfig() error {
	err := c.Role.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) Validate() error {
	err := r.Secret.Validate(venafi.MonitorEngine)
	if err != nil {
		return err
	}

	if r.MaxTTL < r.TTL {
		return fmt.Errorf("max_ttl must be greater than or equal to ttl")
	}

	intermediateCertProvided := r.IntermediateCert != nil
	rootCertProvided := r.RootCert != nil

	if (intermediateCertProvided && rootCertProvided) || (!intermediateCertProvided && !rootCertProvided) {
		return fmt.Errorf("error, must provide exactly one of either the intermediate_certificate or root_certificate blocks: %w", errors.ErrConflictingBlocks)
	}

	return nil
}
