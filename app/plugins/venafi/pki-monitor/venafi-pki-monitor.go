package pki_monitor

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/venafi_wrapper/vcert_wrapper"

	"github.com/opencredo/venafi-vault-wizard/app/github"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func (c *VenafiPKIMonitorConfig) GetDownloadURL() (string, error) {
	var searchSubString string
	if c.BuildArch == "" {
		searchSubString = "linux_optional.zip"
	} else {
		searchSubString = fmt.Sprintf("%s_optional.zip", c.BuildArch)
	}
	return github.GetRelease(
		"Venafi/vault-pki-monitor-venafi",
		c.Version,
		searchSubString,
	)
}

func (c *VenafiPKIMonitorConfig) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	configurePluginSection := report.AddSection("Setting up venafi-pki-monitor")

	venafiClient, err := vcert_wrapper.NewVenafiClient(c.Role.Secret.VenafiSecret, "")
	if err != nil {
		return err
	}

	err = c.Role.Configure(configurePluginSection, c.MountPath, vaultClient, venafiClient)
	if err != nil {
		return err
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
		venafi.MonitorEngine,
	)
	if err != nil {
		return err
	}

	if r.EnforcementPolicy != nil {
		err = ConfigureVenafiPolicy(
			configurePluginSection,
			vaultClient,
			mountPath,
			"default",
			map[string]interface{}{
				"venafi_secret":     r.Secret.Name,
				"zone":              r.EnforcementPolicy.Zone,
				"enforcement_roles": r.Name,
				"defaults_roles":    r.Name,
			},
		)
		if err != nil {
			return err
		}
	}
	if r.ImportPolicy != nil {
		err = ConfigureVenafiPolicy(
			configurePluginSection,
			vaultClient,
			mountPath,
			"visibility",
			map[string]interface{}{
				"venafi_secret": r.Secret.Name,
				"zone":          r.ImportPolicy.Zone,
				"import_roles":  r.Name,
			},
		)
		if err != nil {
			return err
		}
	}

	if r.IntermediateCert != nil {
		err = ConfigureIntermediateCertificate(
			configurePluginSection,
			vaultClient,
			mountPath,
			&r.IntermediateCert.CertificateRequest,
			venafiClient,
			r.IntermediateCert.Zone,
		)
		if err != nil {
			return err
		}
	} else {
		err := ConfigureSelfsignedCertificate(
			configurePluginSection,
			vaultClient,
			mountPath,
			r.RootCert,
		)
		if err != nil {
			return err
		}
	}

	err = ConfigureVenafiRole(
		configurePluginSection,
		vaultClient,
		fmt.Sprintf("%s/roles/%s", mountPath, r.Name),
		r.OptionalConfig.GetAsMap(),
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

		fetchCertSection.Info(fmt.Sprintf("Certificates can be requested using:\nvault write %s common_name=\"test.example.com\"", roleIssuePath))
	}
	return nil
}
