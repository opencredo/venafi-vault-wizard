package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
	mockDL "github.com/opencredo/venafi-vault-wizard/mocks/app/downloader"
	mockPlugin "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
	mockSSH "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/ssh"
)

func TestInstallPlugin(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	vaultSSHClient := new(mockSSH.VaultSSHClient)
	downloader := new(mockDL.PluginDownloader)
	pluginImpl := new(mockPlugin.PluginImpl)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer vaultSSHClient.AssertExpectations(t)
	defer downloader.AssertExpectations(t)
	defer pluginImpl.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginURL = "https://github.com/Venafi/plugin/releases/release.zip"
	var pluginType = "venafi-pki-backend"
	var pluginVersion = "v0.9.0"
	var pluginMountPath = "pki"
	var pluginName = fmt.Sprintf("%s_%s-%s", pluginType, pluginVersion, pluginMountPath)
	var pluginDir = "/etc/plugins"
	var pluginPath = fmt.Sprintf("%s/%s", pluginDir, pluginName)

	pluginImpl.On("GetDownloadURL").Return(pluginURL, pluginVersion, nil)
	downloader.On("DownloadPluginAndUnzip", pluginURL).Return(
		[]byte{0, 1, 2},
		"shashashasha",
		nil,
	)
	vaultSSHClient.On("WriteFile",
		mock.Anything,
		pluginPath,
	).Return(nil)
	vaultSSHClient.On("AddIPCLockCapabilityToFile", pluginPath).Return(nil)
	vaultAPIClient.On("RegisterPlugin", pluginName, pluginName, "shashashasha").Return(nil)
	vaultAPIClient.On("MountPlugin", pluginName, pluginMountPath).Return(nil)

	err := InstallPlugin(&InstallPluginInput{
		VaultClient: vaultAPIClient,
		SSHClients:  []ssh.VaultSSHClient{vaultSSHClient},
		Downloader:  downloader,
		Reporter:    report,
		Plugin: plugins.Plugin{
			Type:      pluginType,
			MountPath: pluginMountPath,
			Impl:      pluginImpl,
		},
		PluginDir:     pluginDir,
		MlockDisabled: false,
	})
	require.NoError(t, err)
}

func reportExpectations(report *mockReport.Report, section *mockReport.Section, check *mockReport.Check) {
	report.On("AddSection", mock.AnythingOfType("string")).Return(section)
	section.On("AddCheck", mock.AnythingOfType("string")).Return(check)
	section.On("Info", mock.AnythingOfType("string")).Maybe()
	check.On("UpdateStatus", mock.AnythingOfType("string")).Maybe()
	check.On("Success", mock.AnythingOfType("string"))
}
