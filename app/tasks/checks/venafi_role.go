package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureVenafiRole(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	rolePath, secretName string,
) error {
	check := reportSection.AddCheck("Adding Venafi role...")

	_, err := vaultClient.WriteValue(rolePath, map[string]interface{}{
		"venafi_secret": secretName,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi role: %s", err))
		return err
	}

	check.Success("Venafi role configured at " + rolePath)
	return nil
}

func VerifyVenafiRole(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	rolePath, secretName string,
) error {
	check := reportSection.AddCheck("Checking Venafi role...")

	data, err := vaultClient.ReadValue(rolePath)
	if err != nil {
		check.Error(fmt.Sprintf("Error retrieving Venafi role: %s", err))
		return err
	}

	if data["venafi_secret"] != secretName {
		check.Error(fmt.Sprintf("The Venafi role's venafi_secret field was not as expected: expected %s got %s", secretName, data["venafi_secret"]))
		return fmt.Errorf("venafi role incorrect")
	}

	check.Success("Venafi role correctly configured at " + rolePath)
	return nil
}