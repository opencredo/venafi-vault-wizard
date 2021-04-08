package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureVenafiSecret(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	secretPath string,
	secretValue map[string]interface{},
) error {
	check := reportSection.AddCheck("Adding Venafi secret...")

	_, err := vaultClient.WriteValue(secretPath, secretValue)
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi secret: %s", err))
		return err
	}

	check.Success("Venafi secret configured at " + secretPath)
	return nil
}

func VerifyVenafiSecret(reportSection reporter.Section, vaultClient api.VaultAPIClient, secretPath, venafiZone string) error {
	check := reportSection.AddCheck("Checking Venafi secret...")

	data, err := vaultClient.ReadValue(secretPath)
	if err != nil {
		check.Error(fmt.Sprintf("Error retrieving Venafi secret: %s", err))
		return err
	}

	if data["zone"] != venafiZone {
		check.Error(fmt.Sprintf("The Venafi secret's zone field is not as expected: expected %s got %s", venafiZone, data["zone"]))
		return fmt.Errorf("venafi secret incorrect")
	}

	check.Success("Venafi secret correctly configured at " + secretPath)
	return nil
}
