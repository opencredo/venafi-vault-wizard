package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureVenafiPolicy(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath, policyName string,
	policyConfig map[string]interface{},
) error {
	check := reportSection.AddCheck("Adding Venafi policy...")

	policyPath := fmt.Sprintf("%s/venafi-policy/%s", mountPath, policyName)
	_, err := vaultClient.WriteValue(policyPath, policyConfig)
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi policy: %s", err))
		return err
	}

	check.Success("Venafi policy configured at " + policyPath)
	return nil
}

func VerifyVenafiPolicy(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath, policyName, secretName, zone string,
) error {
	check := reportSection.AddCheck("Checking Venafi policy...")

	policyPath := fmt.Sprintf("%s/venafi-policy/%s", mountPath, policyName)
	data, err := vaultClient.ReadValue(policyPath)
	if err != nil {
		check.Error(fmt.Sprintf("Error retrieving Venafi policy: %s", err))
		return err
	}

	if data["venafi_secret"] != secretName {
		check.Error(fmt.Sprintf("The Venafi policy's venafi_secret field was not as expected: expected %s got %s", secretName, data["venafi_secret"]))
		return fmt.Errorf("venafi policy incorrect")
	}

	if data["zone"] != zone {
		check.Error(fmt.Sprintf("The Venafi policy's zone field was not as expected: expected %s got %s", zone, data["zone"]))
		return fmt.Errorf("venafi policy incorrect")
	}

	check.Success("Venafi policy correctly configured at " + policyPath)
	return nil
}
