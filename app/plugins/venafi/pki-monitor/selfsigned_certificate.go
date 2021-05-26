package pki_monitor

import (
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureSelfsignedCertificate(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath string,
	request *venafi.CertificateRequest,
) error {
	check := reportSection.AddCheck("Generating self-signed root certificate for the plugin to use...")

	// Generate self-signed root certificate
	_, err := vaultClient.WriteValue(mountPath+"/root/generate/internal", request.ToMap())
	if err != nil {
		check.Errorf("Error generating self-signed root CA: %s", err)
		return err
	}

	check.Success("Root certificate set in Vault")
	return nil
}
