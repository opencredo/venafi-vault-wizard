package tasks

import (
	"fmt"

	"github.com/pterm/pterm"

	"github.com/opencredo/venafi-vault-wizard/helpers/vault"
)

type ConfigureVenafiPKIBackendInput struct {
	VaultClient     vault.Vault
	PluginMountPath string
	SecretName      string
	RoleName        string
	VenafiAPIKey    string
	VenafiZoneID    string
}

func ConfigureVenafiPKIBackend(input *ConfigureVenafiPKIBackendInput) error {
	err := input.addVenafiSecret()
	if err != nil {
		return err
	}

	err = input.addVenafiRole()
	if err != nil {
		return err
	}

	return nil
}

func (i *ConfigureVenafiPKIBackendInput) addVenafiSecret() error {
	spinner, _ := pterm.DefaultSpinner.Start("Adding Venafi secret...")

	secretPath := fmt.Sprintf("%s/venafi/%s", i.PluginMountPath, i.SecretName)
	err := i.VaultClient.WriteValue(secretPath, map[string]interface{}{
		"apikey": i.VenafiAPIKey,
		"zone":   i.VenafiZoneID,
	})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error configuring Venafi secret: %s", err))
		return err
	}

	spinner.Success("Venafi secret configured at " + secretPath)
	return nil
}

func (i *ConfigureVenafiPKIBackendInput) addVenafiRole() error {
	spinner, _ := pterm.DefaultSpinner.Start("Adding Venafi secret...")

	rolePath := fmt.Sprintf("%s/roles/%s", i.PluginMountPath, i.RoleName)
	err := i.VaultClient.WriteValue(rolePath, map[string]interface{}{
		"venafi_secret": i.SecretName,
	})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error configuring Venafi role: %s", err))
		return err
	}

	spinner.Success("Venafi role configured")
	return nil
}
