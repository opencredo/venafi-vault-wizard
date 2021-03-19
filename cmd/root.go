package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func NewRootCommand() *cobra.Command {
	cfg := new(config.GlobalConfig)

	rootCmd := &cobra.Command{
		Use:   "vvw",
		Short: "Venafi Vault Wizard",
		Long:  "VVW is a wizard to automate the installation and verification of Venafi PKI plugins for HashiCorp Vault.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// TODO: bind all the Viper stuff here so the config can get overwritten by viper values
			return nil
		},
	}

	setUpGlobalFlags(rootCmd, cfg)

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Installs a Venafi plugin to Vault",
		Long:  "Installs a plugin to allow Vault to request certificates from Venafi, or to provision them on behalf of Venafi",
	}

	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(
		newInstallPKIBackendCmd(cfg),
		newInstallPKIMonitorCommand(cfg),
	)

	return rootCmd
}

func setUpGlobalFlags(cmd *cobra.Command, config *config.GlobalConfig) {
	flags := cmd.PersistentFlags()

	flags.StringVar(
		&config.VaultAddress,
		"vaultAddress",
		"https://127.0.0.1:8200",
		"Vault HTTP API endpoint",
	)
	flags.StringVar(
		&config.VaultToken,
		"vaultToken",
		"root",
		"Token used to authenticate with Vault",
	)
	flags.StringVar(
		&config.SSHUser,
		"sshUser",
		"username",
		"Username with which to log into Vault server over SSH (must have sudo privileges)",
	)
	flags.StringVar(
		&config.SSHPassword,
		"sshPassword",
		"password",
		"Password for SSH user to log into Vault server with",
	)
	flags.UintVar(
		&config.SSHPort,
		"sshPort",
		22,
		"Port on which SSH is running on the Vault server",
	)
	flags.StringVar(
		&config.VenafiAPIKey,
		"venafiAPIKey",
		"",
		"API Key used to access Venafi Cloud",
	)
	flags.StringVar(
		&config.VenafiZone,
		"venafiZone",
		"",
		"Venafi Cloud Project Zone in which to create certificates",
	)
	flags.StringVar(
		&config.MountPath,
		"vaultMountPath",
		"venafi-pki",
		"Vault path at which to mount the Venafi plugin",
	)
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
