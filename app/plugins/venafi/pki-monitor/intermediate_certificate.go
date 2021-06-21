package pki_monitor

import (
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper"
	"time"

	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func ConfigureIntermediateCertificate(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath string,
	request *venafi.CertificateRequest,
	venafiClient venafi_wrapper.VenafiWrapper,
	zone string,
) error {
	check := reportSection.AddCheck("Generating a CSR for an intermediate certificate from Venafi...")

	// Get intermediate CSR from plugin
	data, err := vaultClient.WriteValue(mountPath+"/intermediate/generate/internal", request.ToMap())
	if err != nil {
		check.Errorf("Error generating subordinate CSR for an intermediate CA: %s", err)
		return err
	}
	pluginCSR := data["csr"].(string)

	// Turn plugin provided CSR into a CSR that Venafi's vcert_wrapper understands
	enrollReq := &certificate.Request{
		CsrOrigin: certificate.UserProvidedCSR,
	}
	err = enrollReq.SetCSR([]byte(pluginCSR))
	if err != nil {
		check.Errorf("Error parsing intermediate CSR provided by Vault: %s", err)
		return err
	}

	check.UpdateStatus("CSR generated, requesting intermediate certificate from Venafi...")

	// Submit request to venafi for the intermediate cert
	requestID, err := venafiClient.RequestCertificate(enrollReq, zone)
	if err != nil {
		check.Errorf("Error requesting intermediate certificate from Venafi: %s", err)
		return err
	}

	// Wait up to 3 minutes for request to complete and get the certificate back
	venafiPEMs, err := venafiClient.RetrieveCertificate(&certificate.Request{
		PickupID: requestID,
		Timeout:  180 * time.Second,
	}, zone)
	if err != nil {
		check.Errorf("Error retrieving intermediate certificate from Venafi: %s", err)
		return err
	}

	_, err = vaultClient.WriteValue(mountPath+"/intermediate/set-signed", map[string]interface{}{
		"certificate": venafiPEMs.Certificate,
	})
	if err != nil {
		check.Errorf("Error setting intermediate certificate in Vault: %s", err)
		return err
	}

	check.Success("Intermediate certificate set in Vault")
	return nil
}

func VerifyIntermediateCertificate(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	policyPath, secretName string,
) error {
	return nil
}
