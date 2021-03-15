package tasks

import (
	"fmt"
	"testing"

	"github.com/pterm/pterm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/mocks"
)

func TestInstallPlugin(t *testing.T) {
	pterm.DisableOutput()

	vaultAPIClient := new(mocks.VaultAPIClient)
	vaultSSHClient := new(mocks.VaultSSHClient)
	downloader := new(mocks.PluginDownloader)
	defer vaultAPIClient.AssertExpectations(t)
	defer vaultSSHClient.AssertExpectations(t)
	defer downloader.AssertExpectations(t)

	var pluginURL = "https://github.com/Venafi/plugin/releases/release.zip"
	var pluginName = "venafi-pki-backend"
	var pluginMountPath = "pki"
	var pluginDir = "/etc/plugins"
	var pluginPath = fmt.Sprintf("%s/%s", pluginDir, pluginName)

	vaultAPIClient.On("GetPluginDir").Return(pluginDir, nil)
	downloader.On("DownloadPluginAndUnzip", pluginURL).Return(
		[]byte{0, 1, 2},
		"abcdefghijk",
		nil,
	)
	vaultSSHClient.On("WriteFile",
		mock.Anything,
		pluginPath,
	).Return(nil)
	vaultAPIClient.On("IsMLockDisabled").Return(false, nil)
	vaultSSHClient.On("AddIPCLockCapabilityToFile", pluginPath).Return(nil)
	vaultAPIClient.On("RegisterPlugin", pluginName, pluginName, "abcdefghijk").Return(nil)
	vaultAPIClient.On("MountPlugin", pluginName, pluginMountPath).Return(nil)

	err := InstallPlugin(&InstallPluginInput{
		VaultClient:     vaultAPIClient,
		SSHClient:       vaultSSHClient,
		Downloader:      downloader,
		PluginURL:       pluginURL,
		PluginName:      pluginName,
		PluginMountPath: pluginMountPath,
	})
	require.NoError(t, err)
}
