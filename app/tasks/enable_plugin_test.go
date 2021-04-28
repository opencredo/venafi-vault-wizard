package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/vault"
	mockPlugin "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
)

func TestEnablePlugin_first_install(t *testing.T) {
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
		Impl:      pluginImpl,
	}
	var sha = "shashashasha"

	// Check for registered plugin, it isn't in catalog for this mount point
	vaultAPIClient.On("GetPlugin", pluginMock.GetCatalogName()).
		Return(nil, vault.ErrNotFound)
	// Should try to register it
	vaultAPIClient.On("RegisterPlugin", pluginMock.GetCatalogName(), pluginMock.GetFileName(), sha).
		Return(nil)

	err := EnablePlugin(&EnablePluginInput{
		VaultClient: vaultAPIClient,
		Reporter:    report,
		Plugin:      pluginMock,
		SHA:         sha,
	})
	require.NoError(t, err)
}

func TestEnablePlugin_already_installed_correct_version(t *testing.T) {
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
		Impl:      pluginImpl,
	}
	var sha = "shashashasha"

	// Check for registered plugin, it's already there
	vaultAPIClient.On("GetPlugin", pluginMock.GetCatalogName()).
		Return(
			map[string]interface{}{
				"command": pluginMock.GetFileName(),
				"sha":     sha,
			},
			nil,
		)

	err := EnablePlugin(&EnablePluginInput{
		VaultClient: vaultAPIClient,
		Reporter:    report,
		Plugin:      pluginMock,
		SHA:         sha,
	})
	require.NoError(t, err)
}

func TestEnablePlugin_already_installed_wrong_version(t *testing.T) {
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
		Impl:      pluginImpl,
	}
	var sha = "shashashasha"

	// Check for registered plugin, it isn't in catalog for this mount point
	vaultAPIClient.On("GetPlugin", pluginMock.GetCatalogName()).
		Return(
			map[string]interface{}{
				"command": pluginMock.Type + "_v0.8.3",
				"sha":     "wrongsha",
			},
			nil,
		)
	// Should try to register it
	vaultAPIClient.On("RegisterPlugin", pluginMock.GetCatalogName(), pluginMock.GetFileName(), sha).
		Return(nil)
	// Then should try to reload it as it's probably already in use
	vaultAPIClient.On("ReloadPlugin", pluginMock.GetCatalogName()).Return(nil)

	err := EnablePlugin(&EnablePluginInput{
		VaultClient: vaultAPIClient,
		Reporter:    report,
		Plugin:      pluginMock,
		SHA:         sha,
	})
	require.NoError(t, err)
}
