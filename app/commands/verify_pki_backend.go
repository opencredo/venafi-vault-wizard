package commands

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
)

func VerifyPKIBackend(configuration *config.Config) {
	report := pretty.NewReport()

	sshClient, vaultClient, closeFunc, err := tasks.GetClients(&configuration.Vault, report)
	if err != nil {
		return
	}
	defer closeFunc()

	err = tasks.VerifyPluginInstalled(&tasks.VerifyPluginInstalledInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Reporter:        report,
		PluginName:      "venafi-pki-backend",
		PluginMountPath: configuration.PKIBackend.MountPath,
	})
	if err != nil {
		return
	}

	err = tasks.CheckVenafiPKIBackend(&tasks.CheckVenafiPKIBackendInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: configuration.PKIBackend.MountPath,
		// TODO: support multiple roles
		SecretName:  configuration.PKIBackend.Roles[0].Secret.Name,
		SecretValue: configuration.PKIBackend.Roles[0].Secret.GetAsMap(),
		RoleName:    configuration.PKIBackend.Roles[0].Name,
	})
	if err != nil {
		return
	}

	err = tasks.FetchVenafiCertificate(&tasks.FetchVenafiCertificateInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: configuration.PKIBackend.MountPath,
		RoleName:        configuration.PKIBackend.Roles[0].Name,
		// TODO: support zero or multiple test certs
		CommonName: configuration.PKIBackend.Roles[0].TestCerts[0].CommonName,
	})
	if err != nil {
		return
	}

	report.Finish(
		fmt.Sprintf(
			"Finished! You can request a certificate using:\n$ vault write %s/issue/%s common_name=\"%s\"\n",
			configuration.PKIBackend.MountPath,
			configuration.PKIBackend.Roles[0].Name,
			configuration.PKIBackend.Roles[0].TestCerts[0].CommonName,
		),
		"Success! Vault is configured to work with Venafi",
	)
}
