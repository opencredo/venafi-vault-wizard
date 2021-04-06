package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

type CheckVenafiPKIBackendInput struct {
	VaultClient     api.VaultAPIClient
	Reporter        reporter.Report
	PluginMountPath string
	SecretName      string
	RoleName        string
	VenafiAPIKey    string
	VenafiZoneID    string
}

func CheckVenafiPKIBackend(input *CheckVenafiPKIBackendInput) error {
	configurePluginSection := input.Reporter.AddSection("Setting up Venafi PKI backend")

	err := checks.VerifyVenafiSecret(configurePluginSection, input.VaultClient, fmt.Sprintf("%s/venafi/%s", input.PluginMountPath, input.SecretName), input.VenafiZoneID)
	if err != nil {
		return err
	}

	err = checks.VerifyVenafiRole(
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
