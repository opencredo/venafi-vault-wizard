package commands

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func Apply(configuration *config.Config) {
	report := pretty.NewReport()

	sshClients, vaultClient, closeFunc, err := tasks.GetClients(&configuration.Vault, report)
	if err != nil {
		return
	}
	defer closeFunc()

	// TODO: try to ascertain whether we have SSH connections to every replica

	checkConfigSection := report.AddSection("Checking Vault server config")
	pluginDir, err := checks.GetPluginDir(checkConfigSection, vaultClient)
	if err != nil {
		return
	}

	mlockDisabledCheck := checkConfigSection.AddCheck("Checking if mlock is disabled...")
	mlockDisabled, err := vaultClient.IsMLockDisabled()
	if err != nil {
		mlockDisabledCheck.Error(fmt.Sprintf("Error checking whether mlock is disabled: %s", err))
		return
	}
	if !mlockDisabled {
		mlockDisabledCheck.Warning("mlock is disabled in the Vault server config, should be enabled for production")
	} else {
		mlockDisabledCheck.Success("mlock is enabled in the Vault server config")
	}

	checkConfigSection.Info(fmt.Sprintf("The Vault server plugin directory is configured as %s\n", pluginDir))

	for _, plugin := range configuration.Plugins {
		pluginName := fmt.Sprintf("%s_v%s", plugin.Type, plugin.Version)
		pluginFileName := fmt.Sprintf("%s-%s", pluginName, plugin.MountPath)
		pluginPath := fmt.Sprintf("%s/%s", pluginDir, pluginFileName)

		checkFilesystemSection := report.AddSection(fmt.Sprintf("Installing plugin %s to Vault server filesystems", plugin.Type))
		checkFilesystemSection.Info(fmt.Sprintf("Plugin filepath is %s\n", pluginPath))

		downloadCheck := checkFilesystemSection.AddCheck("Downloading plugin...")

		pluginURL, err := plugin.Impl.GetDownloadURL()
		if err != nil {
			return
		}
		pluginDownloader := downloader.NewPluginDownloader()
		pluginBytes, sha, err := pluginDownloader.DownloadPluginAndUnzip(pluginURL)
		if err != nil {
			downloadCheck.Error(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
			return
		}

		downloadCheck.Success("Successfully downloaded plugin")

		for i, sshClient := range sshClients {
			installCheck := checkFilesystemSection.AddCheck(fmt.Sprintf("Copying plugin to Vault server %d...", i+1))
			err = sshClient.WriteFile(bytes.NewReader(pluginBytes), pluginPath)
			if err != nil {
				if errors.Is(err, ssh.ErrFileBusy) {
					installCheck.Warning("File already exists and is busy, so cannot overwrite")
				} else {
					installCheck.Error(fmt.Sprintf("Error copying plugin to Vault: %s", err))
					return
				}
			} else {
				installCheck.Success(fmt.Sprintf("Plugin copied to Vault server %d", i+1))
			}

			if !mlockDisabled {
				mlockCheck := checkFilesystemSection.AddCheck(fmt.Sprintf("Mlock is enabled on the Vault server, attempting to add IPC_LOCK capability to plugin on Vault server %d...", i+1))

				err := sshClient.AddIPCLockCapabilityToFile(pluginPath)
				if err != nil {
					mlockCheck.Warning(fmt.Sprintf("Error adding IPC_LOCK capability to plugin on server %d, might cause errors later on: %s", i+1, err))
				} else {
					mlockCheck.Success("IPC_LOCK capability added to plugin")
				}
			}
		}

		enablePluginSection := report.AddSection("Enabling plugin")

		err = checks.InstallPluginInCatalog(enablePluginSection, vaultClient, pluginName, pluginFileName, sha)
		if err != nil {
			return
		}

		err = checks.InstallPluginMount(enablePluginSection, vaultClient, pluginName, plugin.MountPath)
		if err != nil {
			return
		}

		enablePluginSection.Info(fmt.Sprintf("The plugin is enabled in the catalog and mounted at Vault path %s", plugin.MountPath))

		err = plugin.Impl.Configure(report, vaultClient)
		if err != nil {
			return
		}
	}
}
