package tasks

import (
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

type EnablePluginInput struct {
	VaultClient api.VaultAPIClient
	Reporter    reporter.Report
	Plugin      plugins.Plugin
	SHA         string
}

func EnablePlugin(i *EnablePluginInput) error {
	enablePluginSection := i.Reporter.AddSection("Enabling plugin")

	pluginVersionCheck := enablePluginSection.AddCheck("Checking plugin catalog for existing entry...")

	var pluginNeverInstalled = false

	pluginInfo, err := i.VaultClient.GetPlugin(i.Plugin.GetCatalogName())
	if err != nil {
		if !errors.Is(err, vault.ErrNotFound) {
			pluginVersionCheck.Error(fmt.Sprintf("Error checking if plugin is present in catalog: %s", err))
			return err
		}

		pluginVersionCheck.Success("Plugin not yet added to catalog")
		pluginNeverInstalled = true
	} else {
		correctPluginVersionInstalled := (pluginInfo["command"] == i.Plugin.GetFileName()) &&
			(pluginInfo["sha"] == i.SHA)
		if correctPluginVersionInstalled {
			pluginVersionCheck.Success(
				fmt.Sprintf("Version %s of plugin %s already in catalog", i.Plugin.Version, i.Plugin.GetCatalogName()),
			)
			return nil
		}

		pluginVersionCheck.Success(fmt.Sprintf("Plugin command in catalog is currently %s", pluginInfo["command"]))
	}

	err = checks.InstallPluginInCatalog(
		enablePluginSection,
		i.VaultClient,
		i.Plugin.GetCatalogName(),
		i.Plugin.GetFileName(),
		i.SHA,
	)
	if err != nil {
		return err
	}

	if !pluginNeverInstalled {
		err := checks.ReloadPlugin(enablePluginSection, i.VaultClient, i.Plugin.GetCatalogName())
		if err != nil {
			return err
		}
	}

	enablePluginSection.Info(
		fmt.Sprintf("Version %s of plugin %s installed in catalog\n", i.Plugin.Version, i.Plugin.GetCatalogName()),
	)

	return nil
}
