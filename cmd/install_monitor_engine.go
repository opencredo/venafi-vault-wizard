package cmd

import (
	"log"

	"github.com/spf13/cobra"

	appConf "github.com/opencredo/venafi-vault-wizard/app/config"
)

func newInstallPKIMonitorCommand(_ *appConf.GlobalConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "venafi-pki-monitor",
		Short: "Installs venafi-pki-monitor plugin",
		Long:  "Installs the venafi-pki-monitor plugin to allow Vault to generate certificates on behalf of Venafi",
		Run:   installMonitorEngine,
	}
}

func installMonitorEngine(_ *cobra.Command, _ []string) {
	log.Println("unimplemented")
}
