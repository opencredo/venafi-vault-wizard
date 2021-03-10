package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var installPKIMonitorCommand = &cobra.Command{
	Use:   "venafi-pki-monitor",
	Short: "Installs venafi-pki-monitor plugin",
	Long:  "Installs the venafi-pki-monitor plugin to allow Vault to generate certificates on behalf of Venafi",
	Run:   installMonitorEngine,
}

func init() {
	installCommand.AddCommand(installPKIMonitorCommand)
}

func installMonitorEngine(_ *cobra.Command, _ []string) {
	log.Println("unimplemented")
}
