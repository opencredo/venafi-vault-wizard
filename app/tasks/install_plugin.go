package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type InstallPluginToServersInput struct {
	SSHClients    []ssh.VaultSSHClient
	Reporter      reporter.Report
	Plugin        plugins.Plugin
	PluginFile    []byte
	PluginDir     string
	MlockDisabled bool
}

// InstallPluginToServers connects to the Vault servers over SSH and ensures the correct version of the plugin is
// present in the plugin_dir
func InstallPluginToServers(input *InstallPluginToServersInput) error {
	checkFilesystemSection := input.Reporter.AddSection(
		fmt.Sprintf("Installing plugin %s to Vault server filesystems", input.Plugin.Type),
	)

	pluginPath := fmt.Sprintf("%s/%s", input.PluginDir, input.Plugin.GetFileName())

	for i, sshClient := range input.SSHClients {
		err := checks.InstallPluginOnServer(checkFilesystemSection, sshClient, pluginPath, input.PluginFile)
		if err != nil {
			return err
		}

		if !input.MlockDisabled {
			err := checks.InstallPluginMlock(checkFilesystemSection, sshClient, pluginPath)
			if err != nil {
				return err
			}
		}

		checkFilesystemSection.Info(fmt.Sprintf("Plugin copied to Vault server %d at %s\n", i+1, pluginPath))
	}

	return nil
}
