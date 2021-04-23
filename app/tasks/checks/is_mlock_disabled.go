package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func IsMlockDisabled(checkConfigSection reporter.Section, vaultClient api.VaultAPIClient) (bool, error) {
	mlockDisabledCheck := checkConfigSection.AddCheck("Checking if mlock is disabled...")
	mlockDisabled, err := vaultClient.IsMLockDisabled()
	if err != nil {
		mlockDisabledCheck.Error(fmt.Sprintf("Error checking whether mlock is disabled: %s", err))
		return false, err
	}
	if !mlockDisabled {
		mlockDisabledCheck.Warning("mlock is disabled in the Vault server config, should be enabled for production")
	} else {
		mlockDisabledCheck.Success("mlock is enabled in the Vault server config")
	}
	return mlockDisabled, nil
}
