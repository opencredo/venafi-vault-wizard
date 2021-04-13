package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
)

func VerifyPKIBackend(vaultConfig *config.VaultConfig, pluginConfig *config.PKIBackendConfig, venafiConfig config.VenafiConnectionConfig) {
	report := pretty.NewReport()

	sshClient, vaultClient, closeFunc, err := tasks.GetClients(vaultConfig, report)
	if err != nil {
		return
	}
	defer closeFunc()

	err = tasks.VerifyPluginInstalled(&tasks.VerifyPluginInstalledInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Reporter:        report,
		PluginName:      "venafi-pki-backend",
		PluginMountPath: vaultConfig.MountPath,
	})
	if err != nil {
		return
	}

	err = tasks.CheckVenafiPKIBackend(&tasks.CheckVenafiPKIBackendInput{
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

	err = tasks.FetchVenafiCertificate(&tasks.FetchVenafiCertificateInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: vaultConfig.MountPath,
		RoleName:        pluginConfig.RoleName,
		CommonName:      pluginConfig.TestCertificateCommonName,
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
