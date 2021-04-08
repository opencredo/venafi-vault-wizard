package cmd

import (
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func newInstallPKIBackendCmd(vaultConfig *config.VaultConfig) *cobra.Command {
	// Root command for venafi-pki-backend plugin
	pluginConfig := &config.PKIBackendConfig{}
	pkiBackendCmd := &cobra.Command{
		Use:   "venafi-pki-backend",
		Short: "Installs venafi-pki-backend plugin",
		Long:  "Installs the venafi-pki-backend plugin to allow Vault to request certificates from Venafi",
	}
	setUpPKIBackendFlags(pkiBackendCmd, pluginConfig)

	// Variant for Cloud
	cloudConfig := &config.VenafiCloudConfig{}
	cloudCmd := &cobra.Command{
		Use:   "cloud",
		Short: "Installs venafi-pki-backend plugin using Venafi Cloud",
		Long:  "Installs the venafi-pki-backend plugin to allow Vault to request certificates from Venafi Cloud",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := config.ValidateConfigs(vaultConfig, pluginConfig, cloudConfig)
			if err != nil {
				return err
			}

			commands.InstallPKIBackend(vaultConfig, pluginConfig, cloudConfig)
			return nil
		},
	}
	setUpCloudFlags(cloudCmd, cloudConfig)

	// Variant for TPP
	tppConfig := &config.VenafiTPPConfig{}
	tppCmd := &cobra.Command{
		Use:   "tpp",
		Short: "Installs venafi-pki-backend plugin using Venafi TPP",
		Long:  "Installs the venafi-pki-backend plugin to allow Vault to request certificates from Venafi TPP",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := config.ValidateConfigs(vaultConfig, pluginConfig, tppConfig)
			if err != nil {
				return err
			}

			commands.InstallPKIBackend(vaultConfig, pluginConfig, tppConfig)
			return nil
		},
	}
	setUpTPPFlags(tppCmd, tppConfig)

	pkiBackendCmd.AddCommand(cloudCmd)
	pkiBackendCmd.AddCommand(tppCmd)

	return pkiBackendCmd
}

func setUpPKIBackendFlags(cmd *cobra.Command, cfg *config.PKIBackendConfig) {
	flags := cmd.Flags()

	flags.StringVar(
		&cfg.VenafiSecret,
		"venafiSecretName",
		"vvw",
		"Name of Venafi secret to configure in Vault",
	)
	flags.StringVar(
		&cfg.RoleName,
		"roleName",
		"vvw",
		"Name of role to configure in backend, will be used when requesting certificates",
	)
}

func setUpCloudFlags(cmd *cobra.Command, cfg *config.VenafiCloudConfig) {
	flags := cmd.Flags()

	flags.StringVar(
		&cfg.VenafiAPIKey,
		"venafiAPIKey",
		"",
		"API Key used to access Venafi Cloud",
	)
	flags.StringVar(
		&cfg.VenafiZone,
		"venafiZone",
		"",
		"Venafi Cloud Project Zone in which to create certificates",
	)
}

func setUpTPPFlags(cmd *cobra.Command, cfg *config.VenafiTPPConfig) {
	flags := cmd.Flags()

	flags.StringVar(
		&cfg.URL,
		"tppURL",
		"https://tpp.example.com/vedsdk",
		"URL of TPP instance",
	)
	flags.StringVar(
		&cfg.VenafiZone,
		"venafiPolicy",
		"TLS\\Certificates\\HashiCorp Vault",
		"Path to Policy in Venafi TPP with which to create certificates",
	)
	flags.StringVar(
		&cfg.Username,
		"tppUsername",
		"",
		"Username with sufficient permissions to create certificates",
	)
	flags.StringVar(
		&cfg.Password,
		"tppPassword",
		"",
		"Password for user",
	)
}
