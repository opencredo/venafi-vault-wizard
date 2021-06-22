package pki_monitor

import (
	"fmt"
	"github.com/Venafi/vcert/v4/pkg/certificate"
	"github.com/Venafi/vcert/v4/pkg/venafi/tpp"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	mockVenafiWrapper "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins/venafi/venafi_wrapper"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const testCert = "-----BEGIN CERTIFICATE-----\nMIIFQDCCBCigAwIBAgITLwAAAExjVGItPJSAugAAAAAATDANBgkqhkiG9w0BAQsF\nADBNMRMwEQYKCZImiZPyLGQBGRYDY29tMRowGAYKCZImiZPyLGQBGRYKdmVuYWZp\nZGVtbzEaMBgGA1UEAxMRdmVuYWZpZGVtby1UUFAtQ0EwHhcNMjEwNDI5MTA1ODE5\nWhcNMjMwNDI5MTA1ODE5WjAeMRwwGgYDVQQDExN0ZXN0LnZlbmFmaWRlbW8uY29t\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsK68Yp3BpDm/H7EY1uAn\nsu+OFuUBPNKa1XtMf3/Ajx3I8xFFbZOa89kD6i9eHoA+qdP9NeIoOf0UAIXuFnwN\nqfjF1TdbIk3QaoydW09PDv+xyBpLVTCMqSpDAK4ittxOIp3yY1WDAJbqVSCSm/hW\ncMjG6INFXtGcQhvBSL3n2Shm6TjVPmD2FORRFDwe4ax/cyMGy6rwOAEAyUK4n7SC\nLdRIFY9V5EpwjI4bQPGZc/Md2p0wRNQQF6jJt6VjGsWAzV5RsNumBbaMEsgmNOWs\nIWCqW4p7Zq81juVrGabWKeK1QLYOt/XqgYbXFKVkmmfzSUhPakdAdcdOdbpkCZrQ\n9wIDAQABo4ICRjCCAkIwHgYDVR0RBBcwFYITdGVzdC52ZW5hZmlkZW1vLmNvbTAd\nBgNVHQ4EFgQUtzq8zz3NqFExIj3Vgnh6ZcZ3j2wwHwYDVR0jBBgwFoAUg3V6VFgY\nuCIdKHe+7eUpP9ih9f4wgc4GA1UdHwSBxjCBwzCBwKCBvaCBuoaBt2xkYXA6Ly8v\nQ049dmVuYWZpZGVtby1UUFAtQ0EsQ049dHBwLENOPUNEUCxDTj1QdWJsaWMlMjBL\nZXklMjBTZXJ2aWNlcyxDTj1TZXJ2aWNlcyxDTj1Db25maWd1cmF0aW9uLERDPXZl\nbmFmaWRlbW8sREM9Y29tP2NlcnRpZmljYXRlUmV2b2NhdGlvbkxpc3Q/YmFzZT9v\nYmplY3RDbGFzcz1jUkxEaXN0cmlidXRpb25Qb2ludDCBxgYIKwYBBQUHAQEEgbkw\ngbYwgbMGCCsGAQUFBzAChoGmbGRhcDovLy9DTj12ZW5hZmlkZW1vLVRQUC1DQSxD\nTj1BSUEsQ049UHVibGljJTIwS2V5JTIwU2VydmljZXMsQ049U2VydmljZXMsQ049\nQ29uZmlndXJhdGlvbixEQz12ZW5hZmlkZW1vLERDPWNvbT9jQUNlcnRpZmljYXRl\nP2Jhc2U/b2JqZWN0Q2xhc3M9Y2VydGlmaWNhdGlvbkF1dGhvcml0eTAhBgkrBgEE\nAYI3FAIEFB4SAFcAZQBiAFMAZQByAHYAZQByMA4GA1UdDwEB/wQEAwIFoDATBgNV\nHSUEDDAKBggrBgEFBQcDATANBgkqhkiG9w0BAQsFAAOCAQEAA/nWT2AgWgDdnLrC\nTci7fAfo7yxW3QLfoULWUm7k5odQEM80I1aJo+bu+u/dW8ptWkwXUCaiHLgoVmh/\nzCto5GTmMKCNDFvjpgjYUPRItcAptPfstjPsV4jJ8N7oGJ1HYApwdZEy0cC1zKpi\n1i/iZ7iYVeVN+GPF5Sfa/eoCOpha/+8kL4b/hlY1Hpr29oKcurqPsrVLKGHCz55v\neRI58tWIWiG8nzqPK7pCFkw2Vb8DhpeZbjuU1BOcMN4itRereS5dhl/36JBtvdLq\nREed+xyGYi/tPhZ0XMjgL1zIBTAN4nPJKrN2zW/wU8Gh13MuD3HBh9/sE5zeW33D\nVCKviA==\n-----END CERTIFICATE-----"
const testCsr = "-----BEGIN CERTIFICATE REQUEST-----\nMIICtDCCAZwCAQAwbzELMAkGA1UEBhMCR0IxDzANBgNVBAgMBkxvbmRvbjEPMA0G\nA1UEBwwGTG9uZG9uMRIwEAYDVQQKDAlPcGVuQ3JlZG8xDDAKBgNVBAsMA1ZWVzEc\nMBoGA1UEAwwTdGVzdC52ZW5hZmlkZW1vLmNvbTCCASIwDQYJKoZIhvcNAQEBBQAD\nggEPADCCAQoCggEBALgC9BKmIxDRkzIo2o7ofHJ6W9L2zXO9ZtobvxyL4+AEEsxR\nWs8muudbYcjpqSHFqS9oJdPDCeemehyAvB3Tsa2f6BKLJR94uc3lgDD81mDlRCGd\niCvi9vp64IElxlMvCEl4bDKmW7BLIX2lcuUAwcpwDj8s+enu1B4WS/pGGW494+pE\niabZzufxgE3yf+eK6yf7sExXbmVnx9tmOf7NtYm2nXTzaJYCg0Ydhx/Jc7Xp0QY4\nOdyLkVRFggmuNMLkQgY7DTShcuh6kcRsdXN1+De7fxEJDXbYq1lbTEMVLwNoemjz\nGeUyaN1gfsZ4Mzvj768DvjFv73FsLyKkz6KMIp8CAwEAAaAAMA0GCSqGSIb3DQEB\nCwUAA4IBAQAJr+aAredudl+qMtgWeQw6lfHPpVMk3ndt5yu14ff/qaES5PrIsl2x\nrftuQ349N4t+e/J/WwulmdABtIXgj08mw84ryhS/HPMAlh9zrvnLWc5ssI3Zjs8y\nBIfQKOXTrymiUL7xZcJaiMeW2X0Wo/O88S68+uxwSyrMd0nAEffjNWp3bjHQj27b\nza385uUuLszjmxv1MAw/pifttmFg4yKWyWrn0aOLp7LJEWxDa5W+GNMQw/I0L6fC\nBqPKb21V+EAEJfXdpEOloCAAyktWU7mImc74rgZ2gTDMOxHo23t8BXo0k91dHElt\nU0yUK9up1+P29qzAQlxbSRgNyX7TWOKb\n-----END CERTIFICATE REQUEST-----"

func TestConfigureVenafiPKIMonitor(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	venafiClient := new(mockVenafiWrapper.VenafiWrapper)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer venafiClient.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var apiKey = "supersecure API key"
	var pluginMountPath = "pki"
	var requestID = "x509Request"
	var roleName = "roleName"
	var secretName = "pki-monitor"
	var zone = "zone ID"

	var accessToken = "accesstoken"
	var refreshToken = "refreshtoken"

	var url = "http://someurl.com"
	var username = "username"
	var password = "password"

	var intermediateCertGeneratePath = fmt.Sprintf("%s/intermediate/generate/internal", pluginMountPath)
	var intermediateCertSignedPath = fmt.Sprintf("%s/intermediate/set-signed", pluginMountPath)
	var policyDefaultPath = fmt.Sprintf("%s/venafi-policy/default", pluginMountPath)
	var policyVisibilityPath = fmt.Sprintf("%s/venafi-policy/visibility", pluginMountPath)
	var rolePath = fmt.Sprintf("%s/roles/%s", pluginMountPath, roleName)
	var rootCertPath = fmt.Sprintf("%s/root/generate/internal", pluginMountPath)
	var secretPath = fmt.Sprintf("%s/venafi/%s", pluginMountPath, secretName)


	var venafiConnectionConfig = map[string]interface{}{
		"apikey": apiKey,
	}
	var venafiPolicyConfig = map[string]interface{}{
		"defaults_roles": roleName,
		"enforcement_roles": roleName,
		"venafi_secret": secretName,
		"zone": zone,
	}

	var certificateRequest = venafi.CertificateRequest{
		CommonName:   "test.venafidemo.com",
		OU:           "VVW",
		Organisation: "OpenCredo",
		Locality:     "London",
		Province:     "London",
		Country:      "GB",
		TTL:          "1h",
	}

	var tokenResponse = tpp.OauthGetRefreshTokenResponse{
		Access_token:  accessToken,
		Refresh_token: refreshToken,
	}

	var secretMap = map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"url":           url,
		"zone":          zone,
	}

	var certPEMCollection = certificate.PEMCollection{
		Certificate: testCert,
	}

	var importRolesMap = map[string]interface {}{
		"import_roles":  roleName,
		"venafi_secret": secretName,
		"zone":          zone,
	}

	vaultAPIClient.On("WriteValue", secretPath, venafiConnectionConfig).Return(nil, nil)
	vaultAPIClient.On("WriteValue", policyDefaultPath, venafiPolicyConfig).Return(nil, nil)
	venafiClient.On("GetRefreshToken", mock.Anything).Return(tokenResponse, nil)
	vaultAPIClient.On("WriteValue", secretPath, secretMap).Return(nil, nil)
	vaultAPIClient.On("WriteValue", intermediateCertGeneratePath, certificateRequest.ToMap()).Return(map[string]interface{}{"csr": testCsr}, nil)
	vaultAPIClient.On("WriteValue", rootCertPath, certificateRequest.ToMap()).Return(map[string]interface{}{"csr": testCsr}, nil)
	vaultAPIClient.On("WriteValue", rolePath, map[string]interface{}{}).Return(nil, nil)
	venafiClient.On("RequestCertificate", mock.Anything, zone).Return(requestID, nil)
	venafiClient.On("RetrieveCertificate", mock.Anything, zone).Return(&certPEMCollection, nil)
	vaultAPIClient.On("WriteValue", intermediateCertSignedPath, map[string]interface{}{"certificate": testCert}).Return(nil, nil)
	vaultAPIClient.On("WriteValue", policyVisibilityPath, importRolesMap).Return(nil, nil)

	testCases := map[string]struct {
		config VenafiPKIMonitorConfig
	}{
		"pki-monitor cloud enforcement intermediate config": {
			config: VenafiPKIMonitorConfig{
				MountPath: pluginMountPath,
				Role: Role{
					Name: roleName,
					Secret: venafi.VenafiSecret{
						Name: secretName,
						Cloud: &venafi.VenafiCloudConnection{
							APIKey: apiKey,
						},
					},
					EnforcementPolicy: &Policy{
						Zone: zone,
					},
					IntermediateCert: &IntermediateCertRequest{
						Zone: zone,
						CertificateRequest: certificateRequest,
					},
				},
			},
		},
		"pki-monitor cloud enforcement root config": {
			config: VenafiPKIMonitorConfig{
				MountPath: pluginMountPath,
				Role: Role{
					Name: roleName,
					Secret: venafi.VenafiSecret{
						Name: secretName,
						Cloud: &venafi.VenafiCloudConnection{
							APIKey: apiKey,
						},
					},
					EnforcementPolicy: &Policy{
						Zone: zone,
					},
					RootCert: &certificateRequest,
				},
			},
		},
		"pki-monitor cloud enforcement import root config": {
			config: VenafiPKIMonitorConfig{
				MountPath: pluginMountPath,
				Role: Role{
					Name: roleName,
					Secret: venafi.VenafiSecret{
						Name: secretName,
						Cloud: &venafi.VenafiCloudConnection{
							APIKey: apiKey,
						},
					},
					EnforcementPolicy: &Policy{
						Zone: zone,
					},
					ImportPolicy: &Policy{
						Zone: zone,
					},
					RootCert: &certificateRequest,
				},
			},
		},
		"pki-monitor tpp enforcement root config": {
			config: VenafiPKIMonitorConfig{
				MountPath: pluginMountPath,
				Role: Role{
					Name: roleName,
					Secret: venafi.VenafiSecret{
						Name: secretName,
						TPP: &venafi.VenafiTPPConnection{
							URL:      url,
							Username: username,
							Password: password,
							Zone:     zone,
						},
					},
					EnforcementPolicy: &Policy{
						Zone: zone,
					},
					RootCert: &certificateRequest,
				},
			},
		},
	}


	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.config.Role.Configure(section, tc.config.MountPath, vaultAPIClient, venafiClient)
			require.NoError(t, err)
		})
	}
}

func reportExpectations(report *mockReport.Report, section *mockReport.Section, check *mockReport.Check) {
	report.On("AddSection", mock.AnythingOfType("string")).Return(section).Maybe()
	section.On("AddCheck", mock.AnythingOfType("string")).Return(check)
	section.On("Info", mock.AnythingOfType("string")).Maybe()
	check.On("UpdateStatus", mock.AnythingOfType("string")).Maybe()
	check.On("Success", mock.AnythingOfType("string"))
}
