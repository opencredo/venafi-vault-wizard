package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
	mockPlugin "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
	mockSSH "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/ssh"
)

func TestVerifyPluginInstalled(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	vaultSSHClient := new(mockSSH.VaultSSHClient)
	pluginImpl := new(mockPlugin.PluginImpl)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer vaultSSHClient.AssertExpectations(t)
	defer pluginImpl.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginType = "venafi-pki-backend"
	var pluginVersion = "v0.9.0"
	var pluginMountPath = "pki"
	var pluginName = fmt.Sprintf("%s-%s", pluginType, pluginMountPath)
	var pluginFileName = fmt.Sprintf("%s_%s", pluginType, pluginVersion)
	var pluginDir = "/etc/plugins"
	var pluginPath = fmt.Sprintf("%s/%s", pluginDir, pluginFileName)

	vaultSSHClient.On("FileExists", pluginPath).Return(true, nil)
	vaultSSHClient.On("IsIPCLockCapabilityOnFile", pluginPath).Return(true, nil)
	vaultAPIClient.On("GetPlugin", pluginName).Return(
		map[string]interface{}{
			"command": pluginFileName,
		},
		nil,
	)
	vaultAPIClient.On("GetMountPluginName", pluginMountPath).Return(pluginName, nil)

	err := VerifyPluginInstalled(&VerifyPluginInstalledInput{
		VaultClient: vaultAPIClient,
		SSHClients:  []ssh.VaultSSHClient{vaultSSHClient},
		Reporter:    report,
		Plugin: plugins.Plugin{
			Type:      pluginType,
			Version:   pluginVersion,
			MountPath: pluginMountPath,
			Impl:      pluginImpl,
		},
		PluginDir: pluginDir,
	})
	require.NoError(t, err)
}
