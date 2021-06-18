package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type VerifyPluginInstalledInput struct {
	VaultClient   api.VaultAPIClient
	SSHClients    []ssh.VaultSSHClient
	Reporter      reporter.Report
	Plugin        plugins.PluginConfig
	PluginDir     string
	MlockDisabled bool
}

func VerifyPluginInstalled(input *VerifyPluginInstalledInput) error {
	checkFilesystemSection := input.Reporter.AddSection(
		fmt.Sprintf("Checking installation of plugin %s on Vault server filesystem", input.Plugin.Type),
	)

	pluginName := input.Plugin.GetCatalogName()
	pluginFileName := input.Plugin.GetFileName()
	pluginPath := fmt.Sprintf("%s/%s", input.PluginDir, pluginFileName)

	for i, sshClient := range input.SSHClients {
		err := checks.VerifyPluginOnServer(checkFilesystemSection, sshClient, pluginPath)
		if err != nil {
			return err
		}

		if !input.MlockDisabled {
			err := checks.VerifyPluginMlock(checkFilesystemSection, sshClient, pluginPath)
			if err != nil {
				return err
			}
		}

		checkFilesystemSection.Info(fmt.Sprintf("Plugin copied to Vault server %d\n", i+1))
	}

	pluginConfSection := input.Reporter.AddSection("Checking plugin configuration in Vault")

	err := checks.VerifyPluginInCatalog(pluginConfSection, input.VaultClient, pluginName, pluginFileName)
	if err != nil {
		return err
	}

	err = checks.VerifyPluginMount(pluginConfSection, input.VaultClient, pluginName, input.Plugin.MountPath)
	if err != nil {
		return err
	}

	return nil
}
