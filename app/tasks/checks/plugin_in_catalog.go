package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func InstallPluginInCatalog(reportSection reporter.Section, vaultClient api.VaultAPIClient, pluginName, pluginFilename, sha string) error {
	check := reportSection.AddCheck("Enabling plugin in Vault plugin catalog...")
	err := vaultClient.RegisterPlugin(pluginName, pluginFilename, sha)
	if err != nil {
		check.Error(fmt.Sprintf("Error registering plugin in Vault catalog: %s", err))
		return err
	}

	check.Success("Successfully registered plugin in Vault plugin catalog")
	return nil
}

func VerifyPluginInCatalog(reportSection reporter.Section, vaultClient api.VaultAPIClient, pluginName, pluginFilename string) error {
	check := reportSection.AddCheck("Checking whether plugin is enabled in Vault plugin catalog...")
	plugin, err := vaultClient.GetPlugin(pluginName)
	if err != nil {
		check.Error(fmt.Sprintf("Can't look up plugin in Vault plugin catalog: %s", err))
		return err
	}
	if plugin["command"] != pluginFilename {
		check.Error("Plugin enabled, but the currently configured command is incorrect")
		return fmt.Errorf("wrong plugin command configured")
	}
	check.Success("Plugin is enabled in the Vault plugin catalog")
	return nil
}
