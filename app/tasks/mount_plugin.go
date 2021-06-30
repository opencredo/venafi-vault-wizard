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

type MountPluginInput struct {
	VaultClient api.VaultAPIClient
	Reporter    reporter.Report
	Plugin      plugins.PluginConfig
}

func MountPlugin(i *MountPluginInput) error {
	mountPluginSection := i.Reporter.AddSection("Mounting plugin")

	pluginMountCheck := mountPluginSection.AddCheck("Checking if plugin is already mounted...")

	var pluginType string
	if i.Plugin.VaultPluginType == api.PluginTypeSecrets {
		pluginType = i.Plugin.GetCatalogName()
	} else if i.Plugin.VaultPluginType == api.PluginTypeDatabase {
		pluginType = "database"
	}

	pluginName, err := i.VaultClient.GetMountPluginName(i.Plugin.MountPath)
	if err != nil {
		if !errors.Is(err, vault.ErrPluginNotMounted) {
			pluginMountCheck.Errorf("Error checking plugin mount: %s", err)
			return err
		}

		err = checks.InstallPluginMount(
			mountPluginSection,
			i.VaultClient,
			pluginType,
			i.Plugin.MountPath,
		)
		if err != nil {
			return err
		}

		mountPluginSection.Info(fmt.Sprintf("Plugin %s mounted at %s/\n", pluginType, i.Plugin.MountPath))
		return nil
	}

	if pluginName != pluginType {
		pluginMountCheck.Errorf("Mount path %s is using plugin %s", i.Plugin.MountPath, pluginName)
		return vault.ErrMountPathInUse
	}

	pluginMountCheck.Success("Plugin already mounted")
	return nil
}
