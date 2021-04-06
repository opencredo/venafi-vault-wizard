package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/mocks"
)

func TestVerifyPluginInstalled(t *testing.T) {
	vaultAPIClient := new(mocks.VaultAPIClient)
	vaultSSHClient := new(mocks.VaultSSHClient)
	report := new(mocks.Report)
	section := new(mocks.Section)
	check := new(mocks.Check)
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
