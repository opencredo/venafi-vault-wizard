package venafi

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
)

type VenafiSecret struct {
	Name  string                 `hcl:"name,label"`
	Cloud *VenafiCloudConnection `hcl:"venafi_cloud,block"`
	TPP   *VenafiTPPConnection   `hcl:"venafi_tpp,block"`
}

type VenafiConnectionConfig interface {
	GetAsMap() map[string]interface{}
}

func (v *VenafiSecret) Validate() error {
	cloudConnectionProvided := v.Cloud != nil
	tppConnectionProvided := v.TPP != nil

	// Ensure only one of Cloud or TPP is defined
	if (cloudConnectionProvided && tppConnectionProvided) || (!cloudConnectionProvided && !tppConnectionProvided) {
		return fmt.Errorf("error, must provide exactly one of Cloud or TPP connection details: %w", errors.ErrConflictingBlocks)
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

type VenafiCloudConnection struct {
	APIKey string `hcl:"apikey"`
	Zone   string `hcl:"zone"`
}

type VenafiTPPConnection struct {
	URL      string `hcl:"url"`
	Username string `hcl:"username"`
	// TODO: support access token
	Password string `hcl:"password"`
	Policy   string `hcl:"policy"`
}

func (c *VenafiCloudConnection) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("error with Venafi API key: %w", errors.ErrBlankParam)
	}
	if c.Zone == "" {
		return fmt.Errorf("error with Venafi Zone: %w", errors.ErrBlankParam)
	}
	return nil
}

func (c *VenafiTPPConnection) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("error with TPP URL: %w", errors.ErrBlankParam)
	}
	if c.Policy == "" {
		return fmt.Errorf("error with TPP Policy: %w", errors.ErrBlankParam)
	}
	if c.Username == "" {
		return fmt.Errorf("error with TPP Username: %w", errors.ErrBlankParam)
	}
	if c.Password == "" {
		return fmt.Errorf("error with TPP Password: %w", errors.ErrBlankParam)
	}
	return nil
}
