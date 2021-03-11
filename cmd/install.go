package cmd

import "github.com/spf13/cobra"

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Installs a Venafi plugin to Vault",
	Long:  "Installs a plugin to allow Vault to request certificates from Venafi, or to provision them on behalf of Venafi",
}

func init() {
	rootCmd.AddCommand(installCommand)
	installCommand.PersistentFlags().StringVar(
		&sshAddress,
		"sshAddress",
		"vault.local:22",
		"Hostname and port of Vault server",
	)
	installCommand.PersistentFlags().StringVar(
		&venafiAPIKey,
		"venafiAPIKey",
		"",
		"API Key used to access Venafi Cloud",
	)
	installCommand.MarkPersistentFlagRequired("venafiAPIKey")
	installCommand.PersistentFlags().StringVar(
		&venafiZoneID,
		"venafiZone",
		"",
		"Venafi Cloud Project Zone in which to create certificates",
	)
	installCommand.MarkPersistentFlagRequired("venafiZone")
}

var (
	sshAddress   string
	venafiAPIKey string
	venafiZoneID string
)
