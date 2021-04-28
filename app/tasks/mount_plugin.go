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
	Plugin      plugins.Plugin
}

func MountPlugin(i *MountPluginInput) error {
	mountPluginSection := i.Reporter.AddSection("Mounting plugin")

	pluginMountCheck := mountPluginSection.AddCheck("Checking if plugin is already mounted...")

	pluginName, err := i.VaultClient.GetMountPluginName(i.Plugin.MountPath)
	if err != nil {
		if !errors.Is(err, vault.ErrPluginNotMounted) {
			pluginMountCheck.Error(fmt.Sprintf("Error checking plugin mount: %s", err))
			return err
		}

		err = checks.InstallPluginMount(
			mountPluginSection,
			i.VaultClient,
			i.Plugin.GetCatalogName(),
			i.Plugin.MountPath,
		)
		if err != nil {
			return err
		}

		mountPluginSection.Info(fmt.Sprintf("Plugin %s mounted at %s/\n", i.Plugin.GetCatalogName(), i.Plugin.MountPath))
		return nil
	}

	if pluginName != i.Plugin.GetCatalogName() {
		pluginMountCheck.Error(fmt.Sprintf("Mount path %s is using plugin %s", i.Plugin.MountPath, pluginName))
		return vault.ErrMountPathInUse
	}

	pluginMountCheck.Success("Plugin already mounted")
	return nil
}
