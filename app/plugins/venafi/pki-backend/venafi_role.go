package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureVenafiRole(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	rolePath, secretName, zone string,
	optionalParameters map[string]interface{},
) error {
	check := reportSection.AddCheck("Adding Venafi role...")

	optionalParameters["venafi_secret"] = secretName
	optionalParameters["zone"] = zone

	_, err := vaultClient.WriteValue(rolePath, optionalParameters)
	if err != nil {
		check.Errorf("Error configuring Venafi role: %s", err)
		return err
	}

	check.Success("Venafi role configured at " + rolePath)
	return nil
}

func VerifyVenafiRole(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	rolePath, secretName, zone string,
) error {
	check := reportSection.AddCheck("Checking Venafi role...")

	data, err := vaultClient.ReadValue(rolePath)
	if err != nil {
		check.Errorf("Error retrieving Venafi role: %s", err)
		return err
	}

	if data["venafi_secret"] != secretName {
		check.Errorf("The Venafi role's venafi_secret field was not as expected: expected %s got %s", secretName, data["venafi_secret"])
		return fmt.Errorf("venafi role incorrect")
	}

	if data["zone"] != zone {
		check.Errorf("The Venafi role's zone field was not as expected: expected %s got %s", zone, data["zone"])
		return fmt.Errorf("venafi role incorrect")
	}

	check.Success("Venafi role correctly configured at " + rolePath)
	return nil
}
