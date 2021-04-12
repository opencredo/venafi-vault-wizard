package checks

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func RequestVenafiCertificate(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	rolePath, commonName string,
) error {
	check := reportSection.AddCheck("Requesting certificate from Vault...")

	data, err := vaultClient.WriteValue(rolePath, map[string]interface{}{
		"common_name": commonName,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error retrieving certificate from Vault: %s", err))
		return err
	}

	pemBlock, _ := pem.Decode([]byte(data["certificate"].(string)))
	if pemBlock.Type != "CERTIFICATE" {
		check.Error("Expected a certificate to be returned in PEM format")
		return fmt.Errorf("pem type incorrect")
	}
	certificate, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		check.Error(fmt.Sprintf("Error parsing returned certificate: %s", err))
		return err
	}
	if certificate.Subject.CommonName != commonName {
		check.Error(fmt.Sprintf("Certificate's common name was not as expected: expected %s got %s", commonName, certificate.Subject.CommonName))
		return fmt.Errorf("common_name incorrect")
	}

	check.Success("Successfully requested test certificate from Vault")
	return nil
}
