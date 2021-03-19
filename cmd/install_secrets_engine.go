package cmd

import (
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func newInstallPKIBackendCmd(gcfg *config.GlobalConfig) *cobra.Command {
	cfg := &config.PKIBackendConfig{
		GlobalConfig: gcfg,
	}
	cmd := &cobra.Command{
		Use:   "venafi-pki-backend",
		Short: "Installs venafi-pki-backend plugin",
		Long:  "Installs the venafi-pki-backend plugin to allow Vault to request certificates from Venafi",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := cfg.Validate()
			if err != nil {
				return err
			}

			commands.InstallPKIBackend(cfg)
			return nil
		},
	}

	setUpPKIBackendFlags(cmd, cfg)

	return cmd
}

func setUpPKIBackendFlags(cmd *cobra.Command, cfg *config.PKIBackendConfig) {
	flags := cmd.Flags()

	flags.StringVar(
		&cfg.VenafiSecret,
		"venafiSecretName",
		"cloud",
		"Name of Venafi secret to configure in Vault",
	)
	flags.StringVar(
		&cfg.RoleName,
		"roleName",
		"cloud",
		"Name of role to configure in backend, will be used when requesting certificates",
	)
}
