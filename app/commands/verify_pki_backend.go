package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
)

func VerifyPKIBackend(cfg *config.PKIBackendConfig) {
	report := pretty.NewReport()

	sshClient, vaultClient, closeFunc, err := tasks.GetClients(cfg.GlobalConfig, report)
	if err != nil {
		return
	}
	defer closeFunc()

	err = tasks.VerifyPluginInstalled(&tasks.VerifyPluginInstalledInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Reporter:        report,
		PluginName:      "venafi-pki-backend",
		PluginMountPath: cfg.MountPath,
	})
	if err != nil {
		return
	}

	err = tasks.CheckVenafiPKIBackend(&tasks.CheckVenafiPKIBackendInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: cfg.MountPath,
		SecretName:      cfg.VenafiSecret,
		RoleName:        cfg.RoleName,
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
