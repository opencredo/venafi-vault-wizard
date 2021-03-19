package cmd

import (
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func newVerifyPKIMonitorCmd(cfg *config.GlobalConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venafi-pki-monitor",
		Short: "Verifies correct installation of venafi-pki-monitor plugin",
		Long:  "Verifies that the venafi-pki-monitor plugin was correctly configured and that Vault can provision certificates on behalf of Venafi",
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := cfg.Validate()
			if err != nil {
				return err
			}

			commands.VerifyPKIMonitor(cfg)
			return nil
		},
	}

	return cmd
}
