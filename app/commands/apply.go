package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
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
		pluginBytes, sha, err := tasks.DownloadPlugin(&tasks.DownloadPluginInput{
			Reporter: report,
			Plugin:   plugin,
		})
		if err != nil {
			return
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
			return
		}

		err = tasks.EnablePlugin(&tasks.EnablePluginInput{
			VaultClient: vaultClient,
			Reporter:    report,
			Plugin:      plugin,
			SHA:         sha,
		})
		if err != nil {
			return
		}

		err = tasks.MountPlugin(&tasks.MountPluginInput{
			VaultClient: vaultClient,
			Reporter:    report,
			Plugin:      plugin,
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
