package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
	mockSSH "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/ssh"
)

func TestVerifyPluginInstalled(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	vaultSSHClient := new(mockSSH.VaultSSHClient)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer vaultSSHClient.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginName = "venafi-pki-backend"
	var pluginMountPath = "pki"
	var pluginDir = "/etc/plugins"
	var pluginPath = fmt.Sprintf("%s/%s", pluginDir, pluginName)

	vaultAPIClient.On("GetPluginDir").Return(pluginDir, nil)
	vaultSSHClient.On("FileExists", pluginPath).Return(true, nil)
	vaultAPIClient.On("IsMLockDisabled").Return(false, nil)
	vaultSSHClient.On("IsIPCLockCapabilityOnFile", pluginPath).Return(true, nil)
	vaultAPIClient.On("GetPlugin", pluginName).Return(
		map[string]interface{}{
			"command": pluginName,
		},
		nil,
	)
	vaultAPIClient.On("GetMountPluginName", pluginMountPath).Return(pluginName, nil)

	err := VerifyPluginInstalled(&VerifyPluginInstalledInput{
		VaultClient:     vaultAPIClient,
		SSHClient:       vaultSSHClient,
		Reporter:        report,
		PluginName:      pluginName,
		PluginMountPath: pluginMountPath,
	})
	require.NoError(t, err)
}
