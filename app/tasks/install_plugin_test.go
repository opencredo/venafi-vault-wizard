package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
	mockPlugin "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockSSH "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/ssh"
)

func TestInstallPlugin(t *testing.T) {
	vaultSSHClient := new(mockSSH.VaultSSHClient)
	pluginImpl := new(mockPlugin.PluginImpl)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultSSHClient.AssertExpectations(t)
	defer pluginImpl.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginMock = plugins.Plugin{
		Type:      "venafi-pki-backend",
		Version:   "v0.9.0",
		MountPath: "pki",
		Impl:      pluginImpl,
	}
	var pluginDir = "/etc/plugins"
	var pluginPath = fmt.Sprintf("%s/%s", pluginDir, pluginMock.GetFileName())

	vaultSSHClient.On("WriteFile",
		mock.Anything,
		pluginPath,
	).Return(nil)
	vaultSSHClient.On("AddIPCLockCapabilityToFile", pluginPath).Return(nil)

	err := InstallPluginToServers(&InstallPluginToServersInput{
		SSHClients:    []ssh.VaultSSHClient{vaultSSHClient},
		Reporter:      report,
		Plugin:        pluginMock,
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
