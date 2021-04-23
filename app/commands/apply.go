package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
)

func Apply(configuration *config.Config) {
	report := pretty.NewReport()

	sshClients, vaultClient, closeFunc, err := tasks.GetClients(&configuration.Vault, report)
	if err != nil {
		return
	}
	defer closeFunc()
	pluginDownloader := downloader.NewPluginDownloader()

	// TODO: try to ascertain whether we have SSH connections to every replica

	checkConfigSection := report.AddSection("Checking Vault server config")
	pluginDir, err := checks.GetPluginDir(checkConfigSection, vaultClient)
	if err != nil {
		return
	}

	mlockDisabled, err := checks.IsMlockDisabled(checkConfigSection, vaultClient)
	if err != nil {
		return
	}

	checkConfigSection.Info(fmt.Sprintf("The Vault server plugin directory is configured as %s\n", pluginDir))

	for _, plugin := range configuration.Plugins {
		err := tasks.InstallPlugin(&tasks.InstallPluginInput{
			VaultClient:   vaultClient,
			SSHClients:    sshClients,
			Downloader:    pluginDownloader,
			Reporter:      report,
			Plugin:        plugin,
			PluginDir:     pluginDir,
			MlockDisabled: mlockDisabled,
		})
		if err != nil {
			return
		}

		err = tasks.VerifyPluginInstalled(&tasks.VerifyPluginInstalledInput{
			VaultClient:   vaultClient,
			SSHClients:    sshClients,
			Reporter:      report,
			Plugin:        plugin,
			PluginDir:     pluginDir,
			MlockDisabled: mlockDisabled,
		})
		if err != nil {
			return
		}

		err = plugin.Impl.Configure(report, vaultClient)
		if err != nil {
			return
		}

		err = plugin.Impl.Check(report, vaultClient)
		if err != nil {
			return
		}
	}
}
