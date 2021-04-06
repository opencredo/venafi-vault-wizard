package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
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
	checkFilesystemSection := input.Reporter.AddSection("Configuring Vault server filesystem")
	pluginDir, err := checks.GetPluginDir(checkFilesystemSection, input.VaultClient)
	if err != nil {
		return err
	}

	checkFilesystemSection.Info(fmt.Sprintf("The Vault server plugin directory is configured as %s\n", pluginDir))

	// Copy the plugin to the Vault server and add to Vault plugin catalog
	pluginPath := fmt.Sprintf("%s/%s", pluginDir, input.PluginName)
	sha, err := checks.InstallPluginOnServer(
		checkFilesystemSection,
		input.SSHClient,
		input.Downloader,
		pluginPath,
		input.PluginURL,
	)
	if err != nil {
		return err
	}

	err = checks.InstallPluginMlock(checkFilesystemSection, input.VaultClient, input.SSHClient, pluginPath)
	if err != nil {
		return err
	}

	err = checks.InstallPluginInCatalog(checkFilesystemSection, input.VaultClient, input.PluginName, sha)
	if err != nil {
		return err
	}

	checkFilesystemSection.Info(fmt.Sprintf("The Venafi plugin has been installed as %s\n", input.PluginName))

	// Mount backend for plugin in Vault
	pluginConfigureSection := input.Reporter.AddSection("Configuring plugin")
	err = checks.InstallPluginMount(pluginConfigureSection, input.VaultClient, input.PluginName, input.PluginMountPath)
	if err != nil {
		return err
	}

	pluginConfigureSection.Info(fmt.Sprintf("The plugin has been mounted at as %s\n", input.PluginMountPath))

	return nil
}
