package config

import "fmt"

type VenafiPKIBackendConfig struct {
	MountPath string `hcl:"mount_path"`
	Roles     []Role `hcl:"role,block"`
}

type Role struct {
	Name      string               `hcl:"role,label"`
	Secret    VenafiSecret         `hcl:"secret,block"`
	TestCerts []CertificateRequest `hcl:"test_certificate,block"`
}

type CertificateRequest struct {
	CommonName string `hcl:"common_name"`
}

type VenafiSecret struct {
	Name  string                 `hcl:"name,label"`
	Cloud *VenafiCloudConnection `hcl:"venafi_cloud,block"`
	TPP   *VenafiTPPConnection   `hcl:"venafi_tpp,block"`
}

type VenafiConnectionConfig interface {
	GetAsMap() map[string]interface{}
}

func (c *VenafiPKIBackendConfig) Validate() error {
	if c.MountPath == "" {
		return fmt.Errorf("error with plugin mount path: %w", ErrBlankParam)
	}
	if len(c.Roles) == 0 {
		return fmt.Errorf("error at least one role must be provided: %w", ErrBlankParam)
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

func (v *VenafiSecret) Validate() error {
	cloudConnectionProvided := v.Cloud != nil
	tppConnectionProvided := v.TPP != nil

	// Ensure only one of Cloud or TPP is defined
	if (cloudConnectionProvided && tppConnectionProvided) || (!cloudConnectionProvided && !tppConnectionProvided) {
		return fmt.Errorf("error, must provide exactly one of Cloud or TPP connection details: %w", ErrConflictingBlocks)
	}

	if cloudConnectionProvided {
		return v.Cloud.Validate()
	}

	if tppConnectionProvided {
		return v.TPP.Validate()
	}

	return nil
}

func (v *VenafiSecret) GetAsMap() map[string]interface{} {
	if v.Cloud != nil {
		return map[string]interface{}{
			"apikey": v.Cloud.APIKey,
			"zone":   v.Cloud.Zone,
		}
	}

	if v.TPP != nil {
		return map[string]interface{}{
			"url":          v.TPP.URL,
			"zone":         v.TPP.Policy,
			"tpp_user":     v.TPP.Username,
			"tpp_password": v.TPP.Password,
		}
	}

	return nil
}
