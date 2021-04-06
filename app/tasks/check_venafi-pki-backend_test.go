package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/mocks"
)

func TestCheckVenafiPKIBackend(t *testing.T) {
	vaultAPIClient := new(mocks.VaultAPIClient)
	report := new(mocks.Report)
	section := new(mocks.Section)
	check := new(mocks.Check)
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
			"zone":      zone,
			"otherkeys": "extra info",
		}, nil)
	vaultAPIClient.On("ReadValue", rolePath).Return(
		map[string]interface{}{
			"venafi_secret": secretName,
			"otherkeys":     "more info",
		}, nil)

	err := CheckVenafiPKIBackend(&CheckVenafiPKIBackendInput{
		VaultClient:     vaultAPIClient,
		Reporter:        report,
		PluginMountPath: pluginMountPath,
		SecretName:      secretName,
		RoleName:        roleName,
		VenafiZoneID:    zone,
	})
	require.NoError(t, err)
}
