package tasks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/mocks"
)

func TestConfigureVenafiPKIBackend(t *testing.T) {
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
	var venafiConnectionConfig = map[string]interface{}{
		"apikey": "supersecure API key",
		"zone":   "zone ID",
	}

	vaultAPIClient.On("WriteValue", secretPath, venafiConnectionConfig).Return(nil, nil)
	vaultAPIClient.On("WriteValue", rolePath,
		map[string]interface{}{
			"venafi_secret": secretName,
		},
	).Return(nil, nil)

	err := ConfigureVenafiPKIBackend(&ConfigureVenafiPKIBackendInput{
		VaultClient:     vaultAPIClient,
		Reporter:        report,
		PluginMountPath: pluginMountPath,
		SecretName:      secretName,
		SecretValue:     venafiConnectionConfig,
		RoleName:        roleName,
	})
	require.NoError(t, err)
}
