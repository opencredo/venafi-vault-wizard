package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/tasks/checks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

type FetchVenafiCertificateInput struct {
	VaultClient     api.VaultAPIClient
	Reporter        reporter.Report
	PluginMountPath string
	RoleName        string
	CommonName      string
}

func FetchVenafiCertificate(input *FetchVenafiCertificateInput) error {
	fetchCertSection := input.Reporter.AddSection("Fetching test certificate")

	err := checks.RequestVenafiCertificate(
		fetchCertSection,
		input.VaultClient,
		fmt.Sprintf("%s/issue/%s", input.PluginMountPath, input.RoleName),
		input.CommonName,
	)
	if err != nil {
		return err
	}

	return nil
}
