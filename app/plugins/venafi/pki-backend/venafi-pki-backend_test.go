package pki_backend

import (
	"fmt"
	"testing"

	"github.com/Venafi/vcert/v4/pkg/venafi/tpp"
	mockVenafiWrapper "github.com/opencredo/venafi-vault-wizard/mocks/app/plugins/venafi/venafi_wrapper"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockAPI "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/api"
)

func TestConfigureVenafiPKIBackend(t *testing.T) {
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

	var pluginMountPath = "pki"
	var secretName = "pki-backend"
	var secretPath = fmt.Sprintf("%s/venafi/%s", pluginMountPath, secretName)
	var roleName = "roleName"
	var rolePath = fmt.Sprintf("%s/roles/%s", pluginMountPath, roleName)
	var apiKey = "supersecure API key"
	var accessToken = "accesstoken"
	var refreshToken = "refreshtoken"
	var url = "http://someurl.com"
	var username = "username"
	var password = "password"
	var zone = "zone ID"
	var venafiConnectionConfig = map[string]interface{}{
		"apikey": apiKey,
	}

	vaultAPIClient.On("WriteValue", secretPath, venafiConnectionConfig).Return(nil, nil)
	vaultAPIClient.On("WriteValue", rolePath,
		map[string]interface{}{
			"venafi_secret": secretName,
			"zone":          zone,
		},
	).Return(nil, nil)
	venafiClient.On("GetRefreshToken", mock.Anything).Return(tpp.OauthGetRefreshTokenResponse{
		Access_token:  accessToken,
		Refresh_token: refreshToken,
	}, nil)
	vaultAPIClient.On("WriteValue", secretPath,
		map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"url":           url,
			"zone":          zone,
		},
	).Return(nil, nil)

	testCases := map[string]struct {
		config VenafiPKIBackendConfig
	}{
		"pki-backend vaas config": {
			config: VenafiPKIBackendConfig{
				MountPath: pluginMountPath,
				Roles: []Role{
					{
						Name: roleName,
						Zone: zone,
						Secret: venafi.VenafiSecret{
							Name: secretName,
							VaaS: &venafi.VenafiVaaSConnection{
								APIKey: apiKey,
							},
						},
					},
				},
			},
		},
		"pki-backend tpp config": {
			config: VenafiPKIBackendConfig{
				MountPath: pluginMountPath,
				Roles: []Role{
					{
						Name: roleName,
						Zone: zone,
						Secret: venafi.VenafiSecret{
							Name: secretName,
							TPP: &venafi.VenafiTPPConnection{
								URL:      url,
								Username: username,
								Password: password,
								Zone:     zone,
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.config.Roles[0].Configure(section, tc.config.MountPath, vaultAPIClient, venafiClient)
			require.NoError(t, err)
		})
	}
}

func TestCheckVenafiPKIBackend(t *testing.T) {
	vaultAPIClient := new(mockAPI.VaultAPIClient)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultAPIClient.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	var pluginMountPath = "pki"
	var roleName = "roleName"
	var roleIssuePath = fmt.Sprintf("%s/issue/%s", pluginMountPath, roleName)
	var zone = "zone ID"
	var testCSR = venafi.CertificateRequest{
		CommonName:   "test.venafidemo.com",
		OU:           "VVW",
		Organisation: "OpenCredo",
		Locality:     "London",
		Province:     "London",
		Country:      "GB",
		TTL:          "1h",
	}

	vaultAPIClient.On("WriteValue", roleIssuePath, testCSR.ToMap()).
		Return(
			map[string]interface{}{
				"certificate": testCert,
			}, nil,
		)

	config := VenafiPKIBackendConfig{
		MountPath: pluginMountPath,
		Roles: []Role{
			{
				Name:      roleName,
				Zone:      zone,
				TestCerts: []venafi.CertificateRequest{testCSR},
			},
		},
	}
	err := config.Check(report, vaultAPIClient)
	require.NoError(t, err)
}

const testCert = "-----BEGIN CERTIFICATE-----\nMIIFQDCCBCigAwIBAgITLwAAAExjVGItPJSAugAAAAAATDANBgkqhkiG9w0BAQsF\nADBNMRMwEQYKCZImiZPyLGQBGRYDY29tMRowGAYKCZImiZPyLGQBGRYKdmVuYWZp\nZGVtbzEaMBgGA1UEAxMRdmVuYWZpZGVtby1UUFAtQ0EwHhcNMjEwNDI5MTA1ODE5\nWhcNMjMwNDI5MTA1ODE5WjAeMRwwGgYDVQQDExN0ZXN0LnZlbmFmaWRlbW8uY29t\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsK68Yp3BpDm/H7EY1uAn\nsu+OFuUBPNKa1XtMf3/Ajx3I8xFFbZOa89kD6i9eHoA+qdP9NeIoOf0UAIXuFnwN\nqfjF1TdbIk3QaoydW09PDv+xyBpLVTCMqSpDAK4ittxOIp3yY1WDAJbqVSCSm/hW\ncMjG6INFXtGcQhvBSL3n2Shm6TjVPmD2FORRFDwe4ax/cyMGy6rwOAEAyUK4n7SC\nLdRIFY9V5EpwjI4bQPGZc/Md2p0wRNQQF6jJt6VjGsWAzV5RsNumBbaMEsgmNOWs\nIWCqW4p7Zq81juVrGabWKeK1QLYOt/XqgYbXFKVkmmfzSUhPakdAdcdOdbpkCZrQ\n9wIDAQABo4ICRjCCAkIwHgYDVR0RBBcwFYITdGVzdC52ZW5hZmlkZW1vLmNvbTAd\nBgNVHQ4EFgQUtzq8zz3NqFExIj3Vgnh6ZcZ3j2wwHwYDVR0jBBgwFoAUg3V6VFgY\nuCIdKHe+7eUpP9ih9f4wgc4GA1UdHwSBxjCBwzCBwKCBvaCBuoaBt2xkYXA6Ly8v\nQ049dmVuYWZpZGVtby1UUFAtQ0EsQ049dHBwLENOPUNEUCxDTj1QdWJsaWMlMjBL\nZXklMjBTZXJ2aWNlcyxDTj1TZXJ2aWNlcyxDTj1Db25maWd1cmF0aW9uLERDPXZl\nbmFmaWRlbW8sREM9Y29tP2NlcnRpZmljYXRlUmV2b2NhdGlvbkxpc3Q/YmFzZT9v\nYmplY3RDbGFzcz1jUkxEaXN0cmlidXRpb25Qb2ludDCBxgYIKwYBBQUHAQEEgbkw\ngbYwgbMGCCsGAQUFBzAChoGmbGRhcDovLy9DTj12ZW5hZmlkZW1vLVRQUC1DQSxD\nTj1BSUEsQ049UHVibGljJTIwS2V5JTIwU2VydmljZXMsQ049U2VydmljZXMsQ049\nQ29uZmlndXJhdGlvbixEQz12ZW5hZmlkZW1vLERDPWNvbT9jQUNlcnRpZmljYXRl\nP2Jhc2U/b2JqZWN0Q2xhc3M9Y2VydGlmaWNhdGlvbkF1dGhvcml0eTAhBgkrBgEE\nAYI3FAIEFB4SAFcAZQBiAFMAZQByAHYAZQByMA4GA1UdDwEB/wQEAwIFoDATBgNV\nHSUEDDAKBggrBgEFBQcDATANBgkqhkiG9w0BAQsFAAOCAQEAA/nWT2AgWgDdnLrC\nTci7fAfo7yxW3QLfoULWUm7k5odQEM80I1aJo+bu+u/dW8ptWkwXUCaiHLgoVmh/\nzCto5GTmMKCNDFvjpgjYUPRItcAptPfstjPsV4jJ8N7oGJ1HYApwdZEy0cC1zKpi\n1i/iZ7iYVeVN+GPF5Sfa/eoCOpha/+8kL4b/hlY1Hpr29oKcurqPsrVLKGHCz55v\neRI58tWIWiG8nzqPK7pCFkw2Vb8DhpeZbjuU1BOcMN4itRereS5dhl/36JBtvdLq\nREed+xyGYi/tPhZ0XMjgL1zIBTAN4nPJKrN2zW/wU8Gh13MuD3HBh9/sE5zeW33D\nVCKviA==\n-----END CERTIFICATE-----"

func reportExpectations(report *mockReport.Report, section *mockReport.Section, check *mockReport.Check) {
	report.On("AddSection", mock.AnythingOfType("string")).Return(section).Maybe()
	section.On("AddCheck", mock.AnythingOfType("string")).Return(check)
	section.On("Info", mock.AnythingOfType("string")).Maybe()
	check.On("UpdateStatus", mock.AnythingOfType("string")).Maybe()
	check.On("Success", mock.AnythingOfType("string"))
}
