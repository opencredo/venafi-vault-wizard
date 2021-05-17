package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/github"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func (c *VenafiPKIBackendConfig) GetDownloadURL() (string, error) {
	// TODO: allow selecting architectures
	return github.GetRelease(
		"Venafi/vault-pki-backend-venafi",
		c.Version,
		"linux.zip",
	)
}

func (c *VenafiPKIBackendConfig) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	configurePluginSection := report.AddSection("Setting up venafi-pki-backend")

	for _, role := range c.Roles {
		err := venafi.ConfigureVenafiSecret(
			configurePluginSection,
			vaultClient,
			fmt.Sprintf("%s/venafi/%s", c.MountPath, role.Secret.Name),
			role.Secret,
			venafi.SecretsEngine,
		)
		if err != nil {
			return err
		}

		err = ConfigureVenafiRole(
			configurePluginSection,
			vaultClient,
			fmt.Sprintf("%s/roles/%s", c.MountPath, role.Name),
			role.Secret.Name,
			role.Zone,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *VenafiPKIBackendConfig) Check(report reporter.Report, vaultClient api.VaultAPIClient) error {
	for _, role := range c.Roles {
		roleIssuePath := fmt.Sprintf("%s/issue/%s", c.MountPath, role.Name)

		fetchCertSection := report.AddSection(
			fmt.Sprintf("Requesting test certificates from %s", roleIssuePath),
		)
		for _, cert := range role.TestCerts {
			err := venafi.RequestVenafiCertificate(
				fetchCertSection,
				vaultClient,
				roleIssuePath,
				cert,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
