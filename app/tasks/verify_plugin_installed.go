package tasks

import (
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type VerifyPluginInstalledInput struct {
	VaultClient     api.VaultAPIClient
	SSHClient       ssh.VaultSSHClient
	Reporter        reporter.Report
	PluginName      string
	PluginMountPath string
}

func VerifyPluginInstalled(input *VerifyPluginInstalledInput) error {
	checkPluginDirSection := input.Reporter.AddSection("Checking Vault server filesystem")

	pluginDirCheck := checkPluginDirSection.AddCheck("Checking plugin_directory is configured...")
	pluginDir, err := input.VaultClient.GetPluginDir()
	if err != nil {
		if errors.Is(err, vault.ErrPluginDirNotConfigured) {
			pluginDirCheck.Error("The plugin_directory hasn't been configured correctly in the Vault Server Config")
		} else {
			pluginDirCheck.Error(fmt.Sprintf("Error while trying to read plugin_directory: %s", err))
		}
		return err
	}
	pluginDirCheck.Success("Vault Plugin Directory found")

	pluginExistsCheck := checkPluginDirSection.AddCheck("Checking plugin binary exists...")
	pluginPath := fmt.Sprintf("%s/%s", pluginDir, input.PluginName)
	exists, err := input.SSHClient.FileExists(pluginPath)
	if err != nil {
		pluginExistsCheck.Error(fmt.Sprintf("Error checking plugin binary exists: %s", err))
		return err
	}
	if !exists {
		pluginExistsCheck.Error(fmt.Sprintf("Plugin binary does not exist at %s on Vault server", pluginPath))
		return nil
	}
	pluginExistsCheck.Success("Found plugin binary on Vault server")

	ipcLockCapCheck := checkPluginDirSection.AddCheck("Checking whether plugin needs the IPC_LOCK capability and has it...")
	mLockDisabled, err := input.VaultClient.IsMLockDisabled()
	if err != nil {
		ipcLockCapCheck.Error(fmt.Sprintf("Error checking whether mlock is disabled: %s", err))
		return err
	}
	if mLockDisabled {
		ipcLockCapCheck.Warning("Mlock is disabled on the Vault server, should be enabled for production")
	} else {
		capOnFile, err := input.SSHClient.IsIPCLockCapabilityOnFile(pluginPath)
		if err != nil {
			ipcLockCapCheck.Error(fmt.Sprintf("Error checking plugin binary for the IPC_LOCK capabiltiy: %s", err))
			return err
		}

		if capOnFile {
			ipcLockCapCheck.Success("Mlock is enabled and the plugin binary has the IPC_LOCK capability")
		} else {
			ipcLockCapCheck.Warning("Mlock is enabled on Vault server but the plugin does not appear to have the IPC_LOCK capability, however if things seem to work then ignore this warning")
		}
	}

	pluginConfSection := input.Reporter.AddSection("Check Plugin configuration in Vault")
	pluginEnabledCheck := pluginConfSection.AddCheck("Checking whether plugin is enabled in Vault plugin catalog...")
	plugin, err := input.VaultClient.GetPlugin(input.PluginName)
	if err != nil {
		pluginEnabledCheck.Error(fmt.Sprintf("Can't look up plugin in Vault plugin catalog: %s", err))
		return err
	}
	if plugin["command"] != input.PluginName {
		pluginEnabledCheck.Error("Plugin enabled, but the currently configured command is incorrect")
		return nil
	}
	pluginEnabledCheck.Success("Plugin is enabled in the Vault plugin catalog")

	pluginMountCheck := pluginConfSection.AddCheck("Checking plugin is mounted...")
	pluginName, err := input.VaultClient.GetMountPluginName(input.PluginMountPath)
	if err != nil {
		if errors.Is(err, vault.ErrPluginNotMounted) {
			pluginMountCheck.Error(fmt.Sprintf("Plugin is not mounted at %s", input.PluginMountPath))
		} else {
			pluginMountCheck.Error(fmt.Sprintf("Can't check whether the plugin is mounted: %s", err))
		}
		return err
	}
	if pluginName != input.PluginName {
		pluginMountCheck.Error(fmt.Sprintf("Plugin is not mounted at %s", input.PluginMountPath))
		return nil
	}
	pluginMountCheck.Success("Plugin is mounted")

	return nil
}
