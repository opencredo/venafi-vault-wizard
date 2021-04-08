package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
)

const (
	pkiBackendPluginURL = "https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v0.8.3/venafi-pki-backend_v0.8.3_linux.zip"
)

func InstallPKIBackend(vaultConfig *config.VaultConfig, pluginConfig *config.PKIBackendConfig, venafiConfig config.VenafiConnectionConfig) {
	report := pretty.NewReport()

	sshClient, vaultClient, closeFunc, err := tasks.GetClients(vaultConfig, report)
	if err != nil {
		return
	}
	defer closeFunc()

	err = tasks.InstallPlugin(&tasks.InstallPluginInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Downloader:      downloader.NewPluginDownloader(),
		Reporter:        report,
		PluginURL:       pkiBackendPluginURL,
		PluginName:      "venafi-pki-backend",
		PluginMountPath: vaultConfig.MountPath,
	})
	if err != nil {
		return
	}

	err = tasks.ConfigureVenafiPKIBackend(&tasks.ConfigureVenafiPKIBackendInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: vaultConfig.MountPath,
		SecretName:      pluginConfig.VenafiSecret,
		SecretValue:     venafiConfig.GetAsMap(),
		RoleName:        pluginConfig.RoleName,
	})
	if err != nil {
		return
	}

	report.Finish(
		fmt.Sprintf(
			"Finished! You can try and request a certificate using:\n$ vault write %s/issue/%s common_name=\"test.example.com\"\n",
			vaultConfig.MountPath,
			pluginConfig.RoleName,
		),
		"Success! Vault is configured to work with Venafi",
	)
}
