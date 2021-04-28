package checks

import (
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func InstallPluginMount(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	pluginName, pluginMountPath string,
) error {
	check := reportSection.AddCheck("Mounting plugin...")
	err := vaultClient.MountPlugin(pluginName, pluginMountPath)
	if err != nil {
		check.Errorf("Error mounting plugin: %s", err)
		return err
	}

	check.Success("Plugin mounted")
	return nil
}

func VerifyPluginMount(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	pluginName, pluginMountPath string,
) error {
	pluginMountCheck := reportSection.AddCheck("Checking plugin is mounted...")
	actualPluginName, err := vaultClient.GetMountPluginName(pluginMountPath)
	if err != nil {
		if errors.Is(err, vault.ErrPluginNotMounted) {
			pluginMountCheck.Errorf("Plugin is not mounted at %s", pluginMountPath)
		} else {
			pluginMountCheck.Errorf("Can't check whether the plugin is mounted: %s", err)
		}
		return err
	}
	if actualPluginName != pluginName {
		pluginMountCheck.Errorf("Plugin is not mounted at %s", pluginMountPath)
		return fmt.Errorf("wrong plugin mounted")
	}

	pluginMountCheck.Success("Plugin is mounted")
	return nil
}
