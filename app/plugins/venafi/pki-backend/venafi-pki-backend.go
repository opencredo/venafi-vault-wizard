package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/github"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper/vcert_wrapper"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func (c *VenafiPKIBackendConfig) DownloadPlugin() ([]byte, string, error) {
	// TODO: allow selecting architectures
	downloadURL, err := github.GetRelease(
		"Venafi/vault-pki-backend-venafi",
		c.Version,
		"linux.zip",
	)
	if err != nil {
		return nil, "", err
	}

	return downloader.DownloadPluginAndUnzip(downloadURL)
}

func (c *VenafiPKIBackendConfig) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	configurePluginSection := report.AddSection("Setting up venafi-pki-backend")

	for _, role := range c.Roles {
		venafiClient, err := vcert_wrapper.NewVenafiClient(role.Secret, "")
		if err != nil {
			return err
		}

		err = role.Configure(configurePluginSection, c.MountPath, vaultClient, venafiClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Role) Configure(
	configurePluginSection reporter.Section,
	mountPath string,
	vaultClient api.VaultAPIClient,
	venafiClient venafi_wrapper.VenafiWrapper,
) error {
	err := venafi.ConfigureVenafiSecret(
		configurePluginSection,
		vaultClient,
		venafiClient,
		fmt.Sprintf("%s/venafi/%s", mountPath, r.Secret.Name),
		r.Secret,
		venafi.SecretsEngine,
	)
	if err != nil {
		return err
	}

	err = ConfigureVenafiRole(
		configurePluginSection,
		vaultClient,
		fmt.Sprintf("%s/roles/%s", mountPath, r.Name),
		r.Secret.Name,
		r.Zone,
		r.OptionalConfig.GetAsMap(),
	)
	if err != nil {
		return err
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
