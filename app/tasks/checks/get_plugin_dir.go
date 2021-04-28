package checks

import (
	"errors"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func GetPluginDir(reportSection reporter.Section, vaultClient api.VaultAPIClient) (string, error) {
	pluginDirCheck := reportSection.AddCheck("Checking plugin_directory is configured...")
	pluginDir, err := vaultClient.GetPluginDir()
	if err != nil {
		if errors.Is(err, vault.ErrPluginDirNotConfigured) {
			pluginDirCheck.Error("The plugin_directory hasn't been configured correctly in the Vault Server Config")
		} else {
			pluginDirCheck.Errorf("Error while trying to read plugin_directory: %s", err)
		}
		return "", err
	}
	pluginDirCheck.Success("Vault Plugin Directory found")

	return pluginDir, nil
}
