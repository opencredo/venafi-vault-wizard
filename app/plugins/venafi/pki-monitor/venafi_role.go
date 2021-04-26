package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureVenafiRole(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	rolePath string,
	roleParams map[string]interface{},
) error {
	check := reportSection.AddCheck("Adding Venafi role...")

	_, err := vaultClient.WriteValue(rolePath, roleParams)
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
	return nil
}
