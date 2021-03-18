package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

type ConfigureVenafiPKIBackendInput struct {
	VaultClient     api.VaultAPIClient
	Reporter        reporter.Report
	PluginMountPath string
	SecretName      string
	RoleName        string
	VenafiAPIKey    string
	VenafiZoneID    string
}

func ConfigureVenafiPKIBackend(input *ConfigureVenafiPKIBackendInput) error {
	configurePluginSection := input.Reporter.AddSection("Setting up Venafi PKI backend")

	err := input.addVenafiSecret(configurePluginSection)
	if err != nil {
		return err
	}

	err = input.addVenafiRole(configurePluginSection)
	if err != nil {
		return err
	}

	return nil
}

func (i *ConfigureVenafiPKIBackendInput) addVenafiSecret(reportSection reporter.Section) error {
	check := reportSection.AddCheck("Adding Venafi secret...")

	secretPath := fmt.Sprintf("%s/venafi/%s", i.PluginMountPath, i.SecretName)
	_, err := i.VaultClient.WriteValue(secretPath, map[string]interface{}{
		"apikey": i.VenafiAPIKey,
		"zone":   i.VenafiZoneID,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi secret: %s", err))
		return err
	}

	check.Success("Venafi secret configured at " + secretPath)
	return nil
}

func (i *ConfigureVenafiPKIBackendInput) addVenafiRole(reportSection reporter.Section) error {
	check := reportSection.AddCheck("Adding Venafi role...")

	rolePath := fmt.Sprintf("%s/roles/%s", i.PluginMountPath, i.RoleName)
	_, err := i.VaultClient.WriteValue(rolePath, map[string]interface{}{
		"venafi_secret": i.SecretName,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi role: %s", err))
		return err
	}

	check.Success("Venafi role configured")
	return nil
}
