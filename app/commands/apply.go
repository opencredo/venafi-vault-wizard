package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
)

func Apply(configFilePath string, configBytes []byte) error {
	report := pretty.NewReport()
	// TODO: don't return errors if we don't want usage info printed!

	configParser := config.NewConfigParser(configFilePath, configBytes)
	apiAddr, token, err := configParser.GetVaultAPIDetails()
	if err != nil {
		return err
	}

	connectionSection := report.AddSection("Checking connection to Vault")
	vaultClient, err := checks.GetAPIClient(connectionSection, apiAddr, token)
	if err != nil {
		return err
	}

	configParser.SetVaultClient(vaultClient)

	configuration, err := configParser.GetConfig()
	if err != nil {
		return err
	}

	sshClients, closeFunc, err := checks.GetSSHClients(connectionSection, &configuration.Vault)
	if err != nil {
		return err
	}
	defer closeFunc()

	pluginDownloader := downloader.NewPluginDownloader()

	// TODO: try to ascertain whether we have SSH connections to every replica
	checkConfigSection := report.AddSection("Checking Vault server config")
	pluginDir, err := checks.GetPluginDir(checkConfigSection, vaultClient)
	if err != nil {
		return err
	}

	mlockDisabled, err := checks.IsMlockDisabled(checkConfigSection, vaultClient)
	if err != nil {
		return err
	}

	checkConfigSection.Info(fmt.Sprintf("The Vault server plugin directory is configured as %s\n", pluginDir))

	for _, plugin := range configuration.Plugins {
		err = tasks.ResolveBuildArch(&tasks.ResolveBuildArchInput{
			SSHClients:      sshClients,
			PluginBuildArch: plugin.BuildArch,
			Reporter:        report,
		})
		if err != nil {
			return err
		}

		pluginBytes, sha, err := tasks.DownloadPlugin(&tasks.DownloadPluginInput{
			Downloader: pluginDownloader,
			Reporter:   report,
			Plugin:     plugin,
		})
		if err != nil {
			return err
		}

		err = tasks.InstallPluginToServers(&tasks.InstallPluginToServersInput{
			SSHClients:    sshClients,
			Reporter:      report,
			Plugin:        plugin,
			PluginFile:    pluginBytes,
			PluginDir:     pluginDir,
			MlockDisabled: mlockDisabled,
		})
		if err != nil {
			return err
		}

		err = tasks.EnablePlugin(&tasks.EnablePluginInput{
			VaultClient: vaultClient,
			Reporter:    report,
			Plugin:      plugin,
			SHA:         sha,
		})
		if err != nil {
			return err
		}

		err = tasks.MountPlugin(&tasks.MountPluginInput{
			VaultClient: vaultClient,
			Reporter:    report,
			Plugin:      plugin,
		})
		if err != nil {
			return err
		}

		err = plugin.Impl.Configure(report, vaultClient)
		if err != nil {
			return err
		}

		err = plugin.Impl.Check(report, vaultClient)
		if err != nil {
			return err
		}
	}

	return nil
}
