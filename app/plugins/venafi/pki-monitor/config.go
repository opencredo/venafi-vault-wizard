package pki_backend

import (
	"fmt"

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

	IntermediateCert CertificateRequest `hcl:"intermediate_certificate,block"`

	TestCerts []CertificateRequest `hcl:"test_certificate,block"`

	GenerateLease bool   `hcl:"generate_lease,optional"`
	AllowAnyName  bool   `hcl:"allow_any_name,optional"`
	TTL           string `hcl:"ttl,optional"`
	MaxTTL        string `hcl:"max_ttl,optional"`
}

type Policy struct {
	Zone string `hcl:"zone"`
}

type CertificateRequest struct {
	CommonName   string `hcl:"common_name"`
	OU           string `hcl:"ou"`
	Organisation string `hcl:"organisation"`
	Locality     string `hcl:"locality"`
	Province     string `hcl:"province"`
	Country      string `hcl:"country"`
	TTL          string `hcl:"ttl"`
}

func (c *CertificateRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"common_name":  c.CommonName,
		"ou":           c.OU,
		"organization": c.Organisation,
		"locality":     c.Locality,
		"province":     c.Province,
		"country":      c.Country,
		"ttl":          c.TTL,
	}
}

func (c *VenafiPKIMonitorConfig) ValidateConfig() error {
	err := c.Role.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) Validate() error {
	err := r.Secret.Validate()
	if err != nil {
		return err
	}

	if r.MaxTTL < r.TTL {
		return fmt.Errorf("max_ttl must be greater than or equal to ttl")
	}

	return nil
}
