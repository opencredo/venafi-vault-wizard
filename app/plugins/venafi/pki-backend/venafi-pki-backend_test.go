package pki_backend

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
)

func TestConfigureVenafiPKIBackend(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginMountPath = "pki"
	var secretName = "cloud"
	var secretPath = fmt.Sprintf("%s/venafi/%s", pluginMountPath, secretName)
	var roleName = "roleName"
	var rolePath = fmt.Sprintf("%s/roles/%s", pluginMountPath, roleName)
	var apiKey = "supersecure API key"
	var zone = "zone ID"
	var venafiConnectionConfig = map[string]interface{}{
		"apikey": apiKey,
	}

	vaultAPIClient.On("WriteValue", secretPath, venafiConnectionConfig).Return(nil, nil)
	vaultAPIClient.On("WriteValue", rolePath,
		map[string]interface{}{
			"venafi_secret": secretName,
			"zone":          zone,
		},
	).Return(nil, nil)

	config := VenafiPKIBackendConfig{
		MountPath: pluginMountPath,
		Roles: []Role{
			{
				Name: roleName,
				Zone: zone,
				Secret: venafi.VenafiSecret{
					Name: secretName,
					Cloud: &venafi.VenafiCloudConnection{
						APIKey: apiKey,
					},
				},
			},
		},
	}
	err := config.Configure(report, vaultAPIClient)
	require.NoError(t, err)
}

func TestCheckVenafiPKIBackend(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginMountPath = "pki"
	var secretName = "cloud"
	var secretPath = fmt.Sprintf("%s/venafi/%s", pluginMountPath, secretName)
	var roleName = "roleName"
	var rolePath = fmt.Sprintf("%s/roles/%s", pluginMountPath, roleName)
	var zone = "zone ID"

	vaultAPIClient.On("ReadValue", secretPath).Return(
		map[string]interface{}{
			"apikey":    "****",
			"otherkeys": "extra info",
		}, nil)
	vaultAPIClient.On("ReadValue", rolePath).Return(
		map[string]interface{}{
			"venafi_secret": secretName,
			"zone":          zone,
			"otherkeys":     "more info",
		}, nil)

	config := VenafiPKIBackendConfig{
		MountPath: pluginMountPath,
		Roles: []Role{
			{
				Name: roleName,
				Zone: zone,
				Secret: venafi.VenafiSecret{
					Name: secretName,
					Cloud: &venafi.VenafiCloudConnection{
						APIKey: "apikey",
					},
				},
			},
		},
	}
	err := config.Check(report, vaultAPIClient)
	require.NoError(t, err)
}

func reportExpectations(report *mockReport.Report, section *mockReport.Section, check *mockReport.Check) {
	report.On("AddSection", mock.AnythingOfType("string")).Return(section)
	section.On("AddCheck", mock.AnythingOfType("string")).Return(check)
	section.On("Info", mock.AnythingOfType("string")).Maybe()
	check.On("UpdateStatus", mock.AnythingOfType("string")).Maybe()
	check.On("Success", mock.AnythingOfType("string"))
}
