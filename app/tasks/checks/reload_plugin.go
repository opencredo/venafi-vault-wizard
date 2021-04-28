package checks

import (
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ReloadPlugin(reportSection reporter.Section, vaultClient api.VaultAPIClient, pluginName string) error {
	pluginReloadCheck := reportSection.AddCheck("Reloading plugin...")

	err := vaultClient.ReloadPlugin(pluginName)
	if err != nil {
		pluginReloadCheck.Errorf("Error reloading plugin: %s", err)
		return err
	}

	pluginReloadCheck.Success("Plugin reloaded")
	return nil
}
