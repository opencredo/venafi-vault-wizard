package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureVenafiPolicy(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath, secretName string,
	policy Policy,
) error {
	check := reportSection.AddCheck("Adding Venafi policy...")

	policyPath := fmt.Sprintf("%s/venafi-policy/%s", mountPath, policy.Name)
	_, err := vaultClient.WriteValue(policyPath, map[string]interface{}{
		"venafi_secret": secretName,
		"zone":          policy.Zone,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi policy: %s", err))
		return err
	}

	check.Success("Venafi role configured at " + policyPath)
	return nil
}

func VerifyVenafiPolicy(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath, secretName string,
	policy Policy,
) error {
	check := reportSection.AddCheck("Checking Venafi policy...")

	policyPath := fmt.Sprintf("%s/venafi-policy/%s", mountPath, policy.Name)
	data, err := vaultClient.ReadValue(policyPath)
	if err != nil {
		check.Error(fmt.Sprintf("Error retrieving Venafi policy: %s", err))
		return err
	}

	if data["venafi_secret"] != secretName {
		check.Error(fmt.Sprintf("The Venafi policy's venafi_secret field was not as expected: expected %s got %s", secretName, data["venafi_secret"]))
		return fmt.Errorf("venafi policy incorrect")
	}

	if data["zone"] != policy.Zone {
		check.Error(fmt.Sprintf("The Venafi policy's zone field was nto as expected: expected %s got %s", policy.Zone, data["zone"]))
		return fmt.Errorf("venafi policy incorrect")
	}

	check.Success("Venafi policy correctly configured at " + policyPath)
	return nil
}
