package pki_backend

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/github"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func (c *VenafiPKIMonitorConfig) GetDownloadURL() (string, error) {
	return github.GetRelease(
		"Venafi/vault-pki-monitor-venafi",
		c.Version,
		"linux_optional.zip",
	)
}

func (c *VenafiPKIMonitorConfig) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	configurePluginSection := report.AddSection("Setting up venafi-pki-monitor")

	err := venafi.ConfigureVenafiSecret(
		configurePluginSection,
		vaultClient,
		fmt.Sprintf("%s/venafi/%s", c.MountPath, c.Role.Secret.Name),
		c.Role.Secret,
		venafi.MonitorEngine,
	)
	if err != nil {
		return err
	}

	if c.Role.EnforcementPolicy != nil {
		err = ConfigureVenafiPolicy(
			configurePluginSection,
			vaultClient,
			c.MountPath,
			"default",
			map[string]interface{}{
				"venafi_secret":     c.Role.Secret.Name,
				"zone":              c.Role.EnforcementPolicy.Zone,
				"enforcement_roles": c.Role.Name,
				"defaults_roles":    c.Role.Name,
			},
		)
		if err != nil {
			return err
		}
	}
	if c.Role.ImportPolicy != nil {
		err = ConfigureVenafiPolicy(
			configurePluginSection,
			vaultClient,
			c.MountPath,
			"visibility",
			map[string]interface{}{
				"venafi_secret": c.Role.Secret.Name,
				"zone":          c.Role.ImportPolicy.Zone,
				"import_roles":  c.Role.Name,
			},
		)
		if err != nil {
			return err
		}
	}

	if c.Role.IntermediateCert != nil {
		err := ConfigureIntermediateCertificate(
			configurePluginSection,
			vaultClient,
			c.Role.Secret,
			c.MountPath,
			c.Role.IntermediateCert,
			c.Role.EnforcementPolicy.Zone,
		)
		if err != nil {
			return err
		}
	} else {
		err := ConfigureSelfsignedCertificate(
			configurePluginSection,
			vaultClient,
			c.MountPath,
			c.Role.RootCert,
		)
		if err != nil {
			return err
		}
	}

	err = ConfigureVenafiRole(
		configurePluginSection,
		vaultClient,
		fmt.Sprintf("%s/roles/%s", c.MountPath, c.Role.Name),
		map[string]interface{}{
			"ttl":            c.Role.TTL,
			"max_ttl":        c.Role.MaxTTL,
			"allow_any_name": c.Role.AllowAnyName,
			"generate_lease": c.Role.GenerateLease,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *VenafiPKIMonitorConfig) Check(report reporter.Report, vaultClient api.VaultAPIClient) error {
	roleIssuePath := fmt.Sprintf("%s/issue/%s", c.MountPath, c.Role.Name)

	fetchCertSection := report.AddSection(
		fmt.Sprintf("Requesting test certificates from %s", roleIssuePath),
	)
	for _, cert := range c.Role.TestCerts {
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
	return nil
}
