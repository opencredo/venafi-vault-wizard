package tasks

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type InstallPluginInput struct {
	VaultClient     api.VaultAPIClient
	SSHClient       ssh.VaultSSHClient
	Downloader      downloader.PluginDownloader
	Reporter        reporter.Report
	PluginURL       string
	PluginName      string
	PluginMountPath string
}

func InstallPlugin(input *InstallPluginInput) error {
	// Check vault before doing anything
	checkVaultSection := input.Reporter.AddSection("Checking Vault")
	pluginDirCheck := checkVaultSection.AddCheck("Checking Vault Plugin Directory...")

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

	checkVaultSection.Info(fmt.Sprintf("The Vault server plugin directory is configured as %s\n", pluginDir))

	// Copy the plugin to the Vault server and add to Vault plugin catalog
	pluginInstallSection := input.Reporter.AddSection("Installing plugin to Vault")
	pluginPath := fmt.Sprintf("%s/%s", pluginDir, input.PluginName)
	sha, err := input.installPlugin(pluginInstallSection, pluginPath)
	if err != nil {
		return err
	}

	err = input.enablePluginMlock(pluginInstallSection, pluginPath)
	if err != nil {
		return err
	}

	err = input.enablePlugin(pluginInstallSection, sha)
	if err != nil {
		return err
	}

	pluginInstallSection.Info(fmt.Sprintf("The Venafi plugin has been installed as %s\n", input.PluginName))

	// Mount backend for plugin in Vault
	pluginConfigureSection := input.Reporter.AddSection("Configuring plugin")
	err = input.mountPlugin(pluginConfigureSection)
	if err != nil {
		return err
	}

	pluginConfigureSection.Info(fmt.Sprintf("The plugin has been mounted at as %s\n", input.PluginMountPath))

	return nil
}

func (i *InstallPluginInput) installPlugin(reportSection reporter.Section, pluginPath string) (string, error) {
	check := reportSection.AddCheck("Installing plugin to Vault server...")
	plugin, sha, err := i.Downloader.DownloadPluginAndUnzip(i.PluginURL)
	if err != nil {
		check.Error(fmt.Sprintf("Could not download plugin from %s: %s", i.PluginURL, err))
		return "", err
	}

	check.UpdateStatus("Successfully downloaded plugin, copying to Vault server over SSH...")

	err = i.SSHClient.WriteFile(bytes.NewReader(plugin), pluginPath)
	if err != nil {
		if errors.Is(err, ssh.ErrFileBusy) {
			check.Warning("File already exists and is busy, so cannot overwrite")
			return sha, nil
		}

		check.Error(fmt.Sprintf("Error copying plugin to Vault: %s", err))
		return "", err
	}

	check.Success("Plugin copied to Vault server")

	return sha, nil
}

func (i *InstallPluginInput) enablePluginMlock(reportSection reporter.Section, pluginFilepath string) error {
	check := reportSection.AddCheck("Checking if mlock is disabled...")
	mlockDisabled, err := i.VaultClient.IsMLockDisabled()
	if err != nil {
		check.Error(fmt.Sprintf("Error checking whether mlock is disabled: %s", err))
		return err
	}

	if !mlockDisabled {
		check.UpdateStatus("Mlock is enabled on the Vault server, attempting to add IPC_LOCK capability to plugin...")

		err := i.SSHClient.AddIPCLockCapabilityToFile(pluginFilepath)
		if err != nil {
			check.Warning(fmt.Sprintf("Error adding IPC_LOCK capability to plugin, might be needed for mlock: %s", err))
			return nil
		}

		check.Success("IPC_LOCK capability added to plugin")
		return nil
	}

	check.Warning("Mlock is disabled on the Vault server, should be enabled for production")
	return nil
}

func (i *InstallPluginInput) enablePlugin(reportSection reporter.Section, sha string) error {
	check := reportSection.AddCheck("Enabling plugin in Vault plugin catalog...")
	err := i.VaultClient.RegisterPlugin(i.PluginName, i.PluginName, sha)
	if err != nil {
		check.Error(fmt.Sprintf("Error registering plugin in Vault catalog: %s", err))
		return err
	}

	check.Success("Successfully registered plugin in Vault plugin catalog")
	return nil
}

func (i *InstallPluginInput) mountPlugin(reportSection reporter.Section) error {
	check := reportSection.AddCheck("Mounting plugin...")
	err := i.VaultClient.MountPlugin(i.PluginName, i.PluginMountPath)
	if err != nil {
		check.Error(fmt.Sprintf("Error mounting plugin: %s", err))
		return err
	}

	check.Success("Plugin mounted")
	return nil
}
