package tasks

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	mockPlugin "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
)

func TestMountPlugin_already_mounted_correctly(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	pluginImpl := new(mockPlugin.PluginImpl)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer pluginImpl.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginMock = plugins.Plugin{
		Type:      "venafi-pki-monitor",
		MountPath: "venafi-pki",
		Version:   "v0.9.0",
		Impl: pluginImpl,
	}

	vaultAPIClient.On("GetMountPluginName", pluginMock.MountPath).Return(pluginMock.GetCatalogName(), nil)

	err := MountPlugin(&MountPluginInput{
		VaultClient: vaultAPIClient,
		Reporter:    report,
		Plugin:      pluginMock,
	})
	require.NoError(t, err)
}

func TestMountPlugin_already_mounted_wrong_plugin(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	pluginImpl := new(mockPlugin.PluginImpl)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer pluginImpl.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)

	reportExpectations(report, section, check)
	check.On("Error", mock.Anything)

	var pluginMock = plugins.Plugin{
		Type:      "venafi-pki-monitor",
		MountPath: "venafi-pki",
		Version:   "v0.9.0",
		Impl: pluginImpl,
	}

	vaultAPIClient.On("GetMountPluginName", pluginMock.MountPath).Return("some wrong plugin", nil)

	err := MountPlugin(&MountPluginInput{
		VaultClient: vaultAPIClient,
		Reporter:    report,
		Plugin:      pluginMock,
	})
	require.Error(t, err)
}

func TestMountPlugin_first_install(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	pluginImpl := new(mockPlugin.PluginImpl)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer pluginImpl.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginMock = plugins.Plugin{
		Type:      "venafi-pki-monitor",
		MountPath: "venafi-pki",
		Version:   "v0.9.0",
		Impl: pluginImpl,
	}

	vaultAPIClient.On("GetMountPluginName", pluginMock.MountPath).
		Return("", vault.ErrPluginNotMounted)
	vaultAPIClient.On("MountPlugin", pluginMock.GetCatalogName(), pluginMock.MountPath).
		Return(nil)

	err := MountPlugin(&MountPluginInput{
		VaultClient: vaultAPIClient,
		Reporter:    report,
		Plugin:      pluginMock,
	})
	require.NoError(t, err)
}