package cmd

import (
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func newVerifyPKIBackendCmd(vaultConfig *config.VaultConfig) *cobra.Command {
	// Root command for venafi-pki-backend plugin
	pluginConfig := &config.PKIBackendConfig{}
	pkiBackendCmd := &cobra.Command{
		Use:   "venafi-pki-backend",
		Short: "Verifies correct installation of venafi-pki-backend plugin",
		Long:  "Verifies that the venafi-pki-backend plugin was correctly configured and that Vault can request certificates from Venafi",
	}
	setUpPKIBackendFlags(pkiBackendCmd, pluginConfig)

	// Variant for Cloud
	cloudConfig := &config.VenafiCloudConfig{}
	cloudCmd := &cobra.Command{
		Use:   "cloud",
		Short: "Verifies correct installation of venafi-pki-backend plugin with Venafi Cloud",
		Long:  "Verifies that the venafi-pki-backend plugin was correctly configured and that Vault can request certificates from Venafi Cloud",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := config.ValidateConfigs(vaultConfig, pluginConfig, cloudConfig)
			if err != nil {
				return err
			}

			commands.VerifyPKIBackend(vaultConfig, pluginConfig, cloudConfig)
			return nil
		},
	}
	setUpCloudFlags(cloudCmd, cloudConfig)

	// Variant for TPP
	tppConfig := &config.VenafiTPPConfig{}
	tppCmd := &cobra.Command{
		Use:   "tpp",
		Short: "Verifies correct installation of venafi-pki-backend plugin with Venafi TPP",
		Long:  "Verifies that the venafi-pki-backend plugin was correctly configured and that Vault can request certificates from Venafi TPP",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := config.ValidateConfigs(vaultConfig, pluginConfig, tppConfig)
			if err != nil {
				return err
			}

			commands.VerifyPKIBackend(vaultConfig, pluginConfig, tppConfig)
			return nil
		},
	}
	setUpTPPFlags(tppCmd, tppConfig)

	pkiBackendCmd.AddCommand(cloudCmd)
	pkiBackendCmd.AddCommand(tppCmd)

	return pkiBackendCmd
}
