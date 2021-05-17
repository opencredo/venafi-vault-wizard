package venafi

import (
	"fmt"

	"github.com/Venafi/vcert/v4/pkg/endpoint"
	"github.com/Venafi/vcert/v4/pkg/venafi/tpp"
	"github.com/opencredo/venafi-vault-wizard/app/config/errors"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

type PluginType int

const (
	SecretsEngine = iota
	MonitorEngine
)

func ConfigureVenafiSecret(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	secretPath string,
	secretValue VenafiConnectionConfig,
	pluginType PluginType,
) error {
	check := reportSection.AddCheck("Adding Venafi secret...")

	secretParameters, err := secretValue.GetAsMap(pluginType)
	if err != nil {
		check.Errorf("Error getting Venafi secret values: %s", err)
		return err
	}
	_, err = vaultClient.WriteValue(secretPath, secretParameters)
	if err != nil {
		check.Errorf("Error configuring Venafi secret: %s", err)
		return err
	}

	check.Success("Venafi secret configured at " + secretPath)
	return nil
}

func VerifyVenafiSecret(reportSection reporter.Section, vaultClient api.VaultAPIClient, secretPath string, secretValue VenafiConnectionConfig) error {
	check := reportSection.AddCheck("Checking Venafi secret...")

	_, err := vaultClient.ReadValue(secretPath)
	if err != nil {
		check.Errorf("Error retrieving Venafi secret: %s", err)
		return err
	}
	// TODO: check this better, maybe try use auth details to do something with vcert?

	check.Success("Venafi secret correctly configured at " + secretPath)
	return nil
}

type VenafiSecret struct {
	Name  string                 `hcl:"name,label"`
	Cloud *VenafiCloudConnection `hcl:"venafi_cloud,block"`
	TPP   *VenafiTPPConnection   `hcl:"venafi_tpp,block"`
}

type VenafiConnectionConfig interface {
	GetAsMap(pluginType PluginType) (map[string]interface{}, error)
}

func (v *VenafiSecret) Validate(pluginType PluginType) error {
	cloudConnectionProvided := v.Cloud != nil
	tppConnectionProvided := v.TPP != nil

	// Ensure only one of Cloud or TPP is defined
	if (cloudConnectionProvided && tppConnectionProvided) || (!cloudConnectionProvided && !tppConnectionProvided) {
		return fmt.Errorf("error, must provide exactly one of Cloud or TPP connection details: %w", errors.ErrConflictingBlocks)
	}

	if cloudConnectionProvided {
		if pluginType == SecretsEngine && v.Cloud.Zone == "" {
			return fmt.Errorf("error, zone must be specified in secret")
		}
		return v.Cloud.Validate()
	}

	if tppConnectionProvided {
		if pluginType == SecretsEngine && v.TPP.Zone == "" {
			return fmt.Errorf("error, zone must be specified in secret")
		}
		return v.TPP.Validate()
	}

	return nil
}

func (v VenafiSecret) GetAsMap(pluginType PluginType) (map[string]interface{}, error) {
	if v.Cloud != nil {
		m := map[string]interface{}{
			"apikey": v.Cloud.APIKey,
		}
		if v.Cloud.Zone != "" {
			m["zone"] = v.Cloud.Zone
		}
		return m, nil
	}

	if v.TPP != nil {
		m, err := v.TPP.getAccessToken(pluginType)
		if err != nil {
			return nil, err
		}

		if v.TPP.Zone != "" {
			m["zone"] = v.TPP.Zone
		}
		return m, nil
	}

	return nil, nil
}

type VenafiCloudConnection struct {
	APIKey string `hcl:"apikey"`
	Zone   string `hcl:"zone,optional"`
}

type VenafiTPPConnection struct {
	URL      string `hcl:"url"`
	Username string `hcl:"username"`
	Password string `hcl:"password"`
	Zone     string `hcl:"zone,optional"`
}

func (c *VenafiCloudConnection) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("error with Venafi API key: %w", errors.ErrBlankParam)
	}
	return nil
}

func (c *VenafiTPPConnection) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("error with TPP URL: %w", errors.ErrBlankParam)
	}
	if c.Username == "" {
		return fmt.Errorf("error with TPP Username: %w", errors.ErrBlankParam)
	}
	if c.Password == "" {
		return fmt.Errorf("error with TPP Password: %w", errors.ErrBlankParam)
	}
	return nil
}

func (c *VenafiTPPConnection) getAccessToken(pluginType PluginType) (map[string]interface{}, error) {
	tppClient, err := tpp.NewConnector(c.URL, c.Zone, false, nil)
	if err != nil {
		return nil, err
	}

	var scope string
	var clientID string
	if pluginType == MonitorEngine {
		scope = "certificate:manage,discover"
		clientID = "hashicorp-vault-monitor-by-venafi"
	} else if pluginType == SecretsEngine {
		scope = "certificate:manage,revoke"
		clientID = "hashicorp-vault-by-venafi"
	} else {
		return nil, fmt.Errorf("unrecognised plugin type")
	}

	tokens, err := tppClient.GetRefreshToken(&endpoint.Authentication{
		User:     c.Username,
		Password: c.Password,
		Scope:    scope,
		ClientId: clientID,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"url":           c.URL,
		"access_token":  tokens.Access_token,
		"refresh_token": tokens.Refresh_token,
	}, nil
}
