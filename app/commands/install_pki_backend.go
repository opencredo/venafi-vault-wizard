package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
)

const (
	pluginName = "venafi-pki-backend"
	pluginURL  = "https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v0.8.3/venafi-pki-backend_v0.8.3_linux.zip"
)

func InstallPKIBackend(cfg *config.PKIBackendConfig) {
	report := pretty.NewReport()

	sshClient, vaultClient, closeFunc, err := tasks.GetClients(cfg.GlobalConfig, report)
	if err != nil {
		return
	}
	defer closeFunc()

	err = tasks.InstallPlugin(&tasks.InstallPluginInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Downloader:      downloader.NewPluginDownloader(),
		Reporter:        report,
		PluginURL:       pluginURL,
		PluginName:      pluginName,
		PluginMountPath: cfg.MountPath,
	})
	if err != nil {
		return
	}

	err = tasks.ConfigureVenafiPKIBackend(&tasks.ConfigureVenafiPKIBackendInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: cfg.MountPath,
		SecretName:      cfg.VenafiSecret,
		RoleName:        cfg.RoleName,
		VenafiAPIKey:    cfg.VenafiAPIKey,
		VenafiZoneID:    cfg.VenafiZone,
	})
	if err != nil {
		return
	}

	report.Finish(
		fmt.Sprintf(
			"Finished! You can try and request a certificate using:\n$ vault write %s/issue/%s common_name=\"test.example.com\"\n",
			cfg.MountPath,
			cfg.RoleName,
		),
		"Success! Vault is configured to work with Venafi",
	)
}
