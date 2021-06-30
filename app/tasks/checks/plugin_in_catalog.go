package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func InstallPluginInCatalog(reportSection reporter.Section, vaultClient api.VaultAPIClient, pluginName, command, sha string, pluginType api.PluginType) error {
	check := reportSection.AddCheck("Enabling plugin in Vault plugin catalog...")
	err := vaultClient.RegisterPlugin(pluginName, command, sha, pluginType)
	if err != nil {
		check.Errorf("Error registering plugin in Vault catalog: %s", err)
		return err
	}

	check.Success("Successfully registered plugin in Vault plugin catalog")
	return nil
}

func VerifyPluginInCatalog(reportSection reporter.Section, vaultClient api.VaultAPIClient, pluginName, command string, pluginType api.PluginType) error {
	check := reportSection.AddCheck("Checking whether plugin is enabled in Vault plugin catalog...")
	plugin, err := vaultClient.GetPlugin(pluginName, pluginType)
	if err != nil {
		check.Errorf("Can't look up plugin in Vault plugin catalog: %s", err)
		return err
	}
	if plugin["command"] != command {
		check.Error("Plugin enabled, but the currently configured command is incorrect")
		return fmt.Errorf("wrong plugin command configured")
	}
	check.Success("Plugin is enabled in the Vault plugin catalog")
	return nil
}
