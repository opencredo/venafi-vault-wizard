package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/mocks"
)

func TestInstallPlugin(t *testing.T) {
	vaultAPIClient := new(mocks.VaultAPIClient)
	vaultSSHClient := new(mocks.VaultSSHClient)
	downloader := new(mocks.PluginDownloader)
	report := new(mocks.Report)
	section := new(mocks.Section)
	check := new(mocks.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer vaultSSHClient.AssertExpectations(t)
	defer downloader.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

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
		Reporter:        report,
		PluginURL:       pluginURL,
		PluginName:      pluginName,
		PluginMountPath: pluginMountPath,
	})
	require.NoError(t, err)
}

func reportExpectations(report *mocks.Report, section *mocks.Section, check *mocks.Check) {
	report.On("AddSection", mock.AnythingOfType("string")).Return(section)
	section.On("AddCheck", mock.AnythingOfType("string")).Return(check)
	section.On("Info", mock.AnythingOfType("string")).Maybe()
	check.On("UpdateStatus", mock.AnythingOfType("string")).Maybe()
	check.On("Success", mock.AnythingOfType("string"))
}
