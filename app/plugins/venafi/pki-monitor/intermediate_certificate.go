package pki_backend

import (
	"net/http"
	"time"

	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/Venafi/vcert/v4/pkg/endpoint"

	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"

	"github.com/Venafi/vcert/v4"
)

func ConfigureIntermediateCertificate(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	venafiSecret venafi.VenafiSecret,
	mountPath string,
	request *venafi.CertificateRequest,
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

	client, err := getVCertClient(venafiSecret, zone)
	if err != nil {
		check.Errorf("Error connecting to Venafi: %s", err)
		return err
	}

	// Turn plugin provided CSR into a CSR that Venafi's vcert understands
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
	requestID, err := client.RequestCertificate(enrollReq)
	if err != nil {
		check.Errorf("Error requesting intermediate certificate from Venafi: %s", err)
		return err
	}

	// Wait up to 3 minutes for request to complete and get the certificate back
	venafiPEMs, err := client.RetrieveCertificate(&certificate.Request{
		PickupID: requestID,
		Timeout:  180 * time.Second,
	})
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

func getVCertClient(secret venafi.VenafiSecret, zone string) (endpoint.Connector, error) {
	var venafiClient endpoint.Connector
	var err error
	if secret.Cloud != nil {
		venafiClient, err = vcert.NewClient(&vcert.Config{
			ConnectorType: endpoint.ConnectorTypeCloud,
			Credentials: &endpoint.Authentication{
				APIKey: secret.Cloud.APIKey,
			},
			Zone: zone,
			// Specify the DefaultClient otherwise vcert creates its own HTTP Client and for some reason this replaces
			// the TLSClientConfig with a non-nil value it gets from somewhere and breaks things with the following:
			// vcert error: server error: server unavailable: Get "https://api.venafi.cloud/v1/useraccounts": net/http: HTTP/1.x transport connection broken: malformed HTTP response
			Client: http.DefaultClient,
		})
		if err != nil {
			return nil, err
		}
	} else {
		venafiClient, err = vcert.NewClient(&vcert.Config{
			ConnectorType: endpoint.ConnectorTypeTPP,
			BaseUrl:       secret.TPP.URL,
			Credentials: &endpoint.Authentication{
				User:     secret.TPP.Username,
				Password: secret.TPP.Password,
			},
			Zone: zone,
		})
		if err != nil {
			return nil, err
		}
	}

	return venafiClient, nil
}
