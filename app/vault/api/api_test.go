package api

import (
	"testing"

	vaultAPI "github.com/hashicorp/vault/api"
	vaultConsts "github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/mocks"
)

func Test_vault_GetPluginDir(t *testing.T) {
	tests := map[string]struct {
		storedPluginDir string
		wantErr         error
	}{
		"correct_config": {
			storedPluginDir: "/etc/vault.d/plugins",
			wantErr:         nil,
		},
		"plugin_dir_not_configured": {
			storedPluginDir: "",
			wantErr:         vault.ErrPluginDirNotConfigured,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vaultAPIClient := new(mocks.VaultAPIWrapper)
			defer vaultAPIClient.AssertExpectations(t)

			vaultClient := getTestVaultClient(vaultAPIClient)

			expectConfigToBeRead(vaultAPIClient, tc.storedPluginDir, false)

			actualDir, err := vaultClient.GetPluginDir()

			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				require.Equal(t, tc.storedPluginDir, actualDir)
			}
		})
	}
}

func getTestVaultClient(apiClient *mocks.VaultAPIWrapper) VaultAPIClient {
	apiClient.On("SetAddress", "apiaddr").Return(nil)
	apiClient.On("SetToken", "tok").Return(nil)

	vaultClient := NewClient(
		&Config{
			APIAddress: "apiaddr",
			Token:      "tok",
		},
		apiClient,
	)

	return vaultClient
}

func expectConfigToBeRead(apiClient *mocks.VaultAPIWrapper, pluginDir string, mlockDisabled bool) {
	apiClient.On("Read", "sys/config/state/sanitized").Return(
		map[string]interface{}{
			"plugin_directory": pluginDir,
			"disable_mlock":    mlockDisabled,
		},
		nil,
	)
}

func Test_vault_RegisterPlugin(t *testing.T) {
	vaultAPIClient := new(mocks.VaultAPIWrapper)
	defer vaultAPIClient.AssertExpectations(t)

	vaultClient := getTestVaultClient(vaultAPIClient)

	vaultAPIClient.On(
		"RegisterPlugin",
		mock.MatchedBy(func(input *vaultAPI.RegisterPluginInput) bool {
			return input.Type == vaultConsts.PluginTypeSecrets
		}),
	).Return(nil)

	err := vaultClient.RegisterPlugin("name", "command", "sha")

	require.NoError(t, err)
}

func Test_vault_MountPlugin(t *testing.T) {
	vaultAPIClient := new(mocks.VaultAPIWrapper)
	defer vaultAPIClient.AssertExpectations(t)

	vaultClient := getTestVaultClient(vaultAPIClient)

	var backendName = "backend"
	var mountPath = "path"

	vaultAPIClient.On(
		"Mount",
		mountPath,
		mock.MatchedBy(func(input *vaultAPI.MountInput) bool {
			return input.Type == backendName
		}),
	).Return(nil)

	err := vaultClient.MountPlugin(backendName, mountPath)
	require.NoError(t, err)
}

func Test_vault_IsMLockDisabled(t *testing.T) {
	vaultAPIClient := new(mocks.VaultAPIWrapper)
	defer vaultAPIClient.AssertExpectations(t)

	vaultClient := getTestVaultClient(vaultAPIClient)

	expectConfigToBeRead(vaultAPIClient, "", false)

	mLockDisabled, err := vaultClient.IsMLockDisabled()
	require.NoError(t, err)

	require.False(t, mLockDisabled)
}
