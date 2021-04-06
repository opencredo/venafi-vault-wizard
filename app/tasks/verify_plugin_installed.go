package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
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
	checkFileSystemSection := input.Reporter.AddSection("Checking Vault server filesystem")
	pluginDir, err := checks.GetPluginDir(checkFileSystemSection, input.VaultClient)
	if err != nil {
		return err
	}

	pluginPath := fmt.Sprintf("%s/%s", pluginDir, input.PluginName)
	err = checks.VerifyPluginOnServer(checkFileSystemSection, input.SSHClient, pluginPath)
	if err != nil {
		return err
	}

	err = checks.VerifyPluginMlock(checkFileSystemSection, input.VaultClient, input.SSHClient, pluginPath)
	if err != nil {
		return err
	}

	pluginConfSection := input.Reporter.AddSection("Check Plugin configuration in Vault")
	err = checks.VerifyPluginInCatalog(pluginConfSection, input.VaultClient, input.PluginName)
	if err != nil {
		return err
	}

	err = checks.VerifyPluginMount(pluginConfSection, input.VaultClient, input.PluginName, input.PluginMountPath)
	if err != nil {
		return err
	}

	return nil
}
