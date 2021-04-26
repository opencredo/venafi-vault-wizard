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
		"linux_strict.zip",
	)
}

func (c *VenafiPKIMonitorConfig) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	configurePluginSection := report.AddSection("Setting up venafi-pki-monitor")

	err := venafi.ConfigureVenafiSecret(
		configurePluginSection,
		vaultClient,
		fmt.Sprintf("%s/venafi/%s", c.MountPath, c.Role.Secret.Name),
		c.Role.Secret,
	)
	if err != nil {
		return err
	}

	err = ConfigureVenafiPolicy(
		configurePluginSection,
		vaultClient,
		c.MountPath,
		c.Role.Secret.Name,
		c.Role.EnforcementPolicy,
	)
	if err != nil {
		return err
	}

	err = ConfigureIntermediateCertificate(
		configurePluginSection,
		vaultClient,
		c.MountPath,
		&c.Role.IntermediateCert,
	)
	if err != nil {
		return err
	}

	if c.Role.ImportPolicy != nil {
		err = ConfigureVenafiPolicy(
			configurePluginSection,
			vaultClient,
			c.MountPath,
			c.Role.Secret.Name,
			*c.Role.ImportPolicy,
		)
		if err != nil {
			return err
		}
	}

	if c.Role.DefaultPolicy != nil {
		err = ConfigureVenafiPolicy(
			configurePluginSection,
			vaultClient,
			c.MountPath,
			c.Role.Secret.Name,
			*c.Role.DefaultPolicy,
		)
		if err != nil {
			return err
		}
	}

	err = ConfigureVenafiRole(
		configurePluginSection,
		vaultClient,
		fmt.Sprintf("%s/role/%s", c.MountPath, c.Role.Name),
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
	panic("implement me")
}
