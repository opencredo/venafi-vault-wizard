package pki_backend

import (
	"fmt"
	"time"

	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/Venafi/vcert/v4/pkg/endpoint"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"

	"github.com/Venafi/vcert/v4"
)

func ConfigureIntermediateCertificate(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	mountPath string,
	request *CertificateRequest,
) error {
	check := reportSection.AddCheck("Generating a CSR for an intermediate certificate from Venafi...")

	// Get intermediate CSR from plugin
	data, err := vaultClient.WriteValue(mountPath+"/intermediate/generate/internal", request.ToMap())
	if err != nil {
		check.Error(fmt.Sprintf("Error configuring Venafi policy: %s", err))
		return err
	}
	pluginCSR := data["csr"].([]byte)

	// Get Venafi client
	client, err := vcert.NewClient(&vcert.Config{
		ConnectorType: endpoint.ConnectorTypeFake,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error connecting to Venafi: %s", err))
		return err
	}

	// Turn plugin provided CSR into a CSR that Venafi's vcert understands
	enrollReq := &certificate.Request{
		CsrOrigin: certificate.UserProvidedCSR,
	}
	err = enrollReq.SetCSR(pluginCSR)
	if err != nil {
		check.Error(fmt.Sprintf("Error parsing intermediate CSR provided by Vault: %s", err))
		return err
	}

	check.UpdateStatus("CSR generated, requesting intermediate certificate from Venafi...")

	// Submit request to venafi for the intermediate cert
	requestID, err := client.RequestCertificate(enrollReq)
	if err != nil {
		check.Error(fmt.Sprintf("Error requesting intermediate certificate from Venafi: %s", err))
		return err
	}

	// Wait up to 3 minutes for request to complete and get the certificate back
	venafiPEMs, err := client.RetrieveCertificate(&certificate.Request{
		PickupID: requestID,
		Timeout:  180 * time.Second,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error retrieving intermediate certificate from Venafi: %s", err))
		return err
	}

	_, err = vaultClient.WriteValue(mountPath+"intermediate/set-signed", map[string]interface{}{
		"certificate": venafiPEMs.Certificate,
	})
	if err != nil {
		check.Error(fmt.Sprintf("Error setting intermediate certificate in Vault: %s", err))
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
