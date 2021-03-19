package cmd

import (
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func newVerifyPKIBackendCmd(gcfg *config.GlobalConfig) *cobra.Command {
	cfg := &config.PKIBackendConfig{
		GlobalConfig: gcfg,
	}

	cmd := &cobra.Command{
		Use:   "venafi-pki-backend",
		Short: "Verifies correct installation of venafi-pki-backend plugin",
		Long:  "Verifies that the venafi-pki-backend plugin was correctly configured and that Vault can request certificates from Venafi",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := cfg.Validate()
			if err != nil {
				return err
			}

			commands.VerifyPKIBackend(cfg)
			return nil
		},
	}

	setUpPKIBackendFlags(cmd, cfg)

	return cmd
}
