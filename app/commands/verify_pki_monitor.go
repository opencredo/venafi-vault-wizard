package commands

import (
	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
)

func VerifyPKIMonitor(cfg *config.VaultConfig) {
	report := pretty.NewReport()

	sshClient, vaultClient, closeFunc, err := tasks.GetClients(cfg, report)
	if err != nil {
		return
	}
	defer closeFunc()

	err = tasks.VerifyPluginInstalled(&tasks.VerifyPluginInstalledInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Reporter:        report,
		PluginName:      "venafi-pki-monitor",
		PluginMountPath: cfg.MountPath,
	})
	if err != nil {
		return
	}
}
