package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

type ConfigureVenafiPKIBackendInput struct {
	VaultClient     api.VaultAPIClient
	Reporter        reporter.Report
	PluginMountPath string
	SecretName      string
	SecretValue     map[string]interface{}
	RoleName        string
}

func ConfigureVenafiPKIBackend(input *ConfigureVenafiPKIBackendInput) error {
	configurePluginSection := input.Reporter.AddSection("Setting up Venafi PKI backend")

	err := checks.ConfigureVenafiSecret(
		configurePluginSection,
		input.VaultClient,
		fmt.Sprintf("%s/venafi/%s", input.PluginMountPath, input.SecretName),
		input.SecretValue,
	)
	if err != nil {
		return err
	}

	err = checks.ConfigureVenafiRole(
		configurePluginSection,
		input.VaultClient,
		fmt.Sprintf("%s/roles/%s", input.PluginMountPath, input.RoleName),
		input.SecretName,
	)
	if err != nil {
		return err
	}

	return nil
}
