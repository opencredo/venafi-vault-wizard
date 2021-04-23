package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type InstallPluginInput struct {
	VaultClient   api.VaultAPIClient
	SSHClients    []ssh.VaultSSHClient
	Downloader    downloader.PluginDownloader
	Reporter      reporter.Report
	Plugin        plugins.Plugin
	PluginDir     string
	MlockDisabled bool
}

func InstallPlugin(input *InstallPluginInput) error {
	checkFilesystemSection := input.Reporter.AddSection(
		fmt.Sprintf("Installing plugin %s to Vault server filesystems", input.Plugin.Type),
	)

	downloadCheck := checkFilesystemSection.AddCheck("Downloading plugin...")

	pluginURL, version, err := input.Plugin.Impl.GetDownloadURL()
	if err != nil {
		return err
	}
	pluginBytes, sha, err := input.Downloader.DownloadPluginAndUnzip(pluginURL)
	if err != nil {
		downloadCheck.Error(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
		return err
	}

	downloadCheck.Success("Successfully downloaded plugin")

	pluginName := fmt.Sprintf("%s_%s-%s", input.Plugin.Type, version, input.Plugin.MountPath)
	pluginPath := fmt.Sprintf("%s/%s", input.PluginDir, pluginName)

	checkFilesystemSection.Info(fmt.Sprintf("Plugin filepath is %s\n", pluginPath))

	for i, sshClient := range input.SSHClients {
		err := checks.InstallPluginOnServer(checkFilesystemSection, sshClient, pluginPath, pluginBytes)
		if err != nil {
			return err
		}

		if !input.MlockDisabled {
			err := checks.InstallPluginMlock(checkFilesystemSection, sshClient, pluginPath)
			if err != nil {
				return err
			}
		}

		checkFilesystemSection.Info(fmt.Sprintf("Plugin copied to Vault server %d\n", i+1))
	}

	enablePluginSection := input.Reporter.AddSection("Enabling plugin")

	err = checks.InstallPluginInCatalog(
		enablePluginSection,
		input.VaultClient,
		pluginName,
		sha,
	)
	if err != nil {
		return err
	}

	err = checks.InstallPluginMount(
		enablePluginSection,
		input.VaultClient,
		pluginName,
		input.Plugin.MountPath,
	)
	if err != nil {
		return err
	}

	enablePluginSection.Info(fmt.Sprintf("The plugin is enabled in the catalog and mounted at Vault path %s", input.Plugin.MountPath))

	return nil
}
