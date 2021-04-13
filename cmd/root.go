package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
	"github.com/opencredo/venafi-vault-wizard/app/config"
)

func NewRootCommand() *cobra.Command {
	var configFile string
	configuration := new(config.Config)

	rootCmd := &cobra.Command{
		Use:   "vvw",
		Short: "Venafi Vault Wizard",
		Long:  "VVW is a wizard to automate the installation and verification of Venafi PKI plugins for HashiCorp Vault.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Parse provided config file
			conf, err := config.NewConfigFromFile(configFile)
			if err != nil {
				return err
			}

			configuration = conf
			return nil
		},
	}

	setUpGlobalFlags(rootCmd, &configFile)

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Installs a Venafi plugin to Vault",
		Long:  "Installs a plugin to allow Vault to request certificates from Venafi, or to provision them on behalf of Venafi",
		Run: func(_ *cobra.Command, _ []string) {
			if configuration.PKIBackend != nil {
				commands.InstallPKIBackend(configuration)
			}
			// TODO: call install PKI monitor when implemented
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verifies correct installation of a Venafi Vault plugin",
		Long:  "Verifies that the installation of either of the Venafi Vault plugins was successful and that it is configured correctly",
		Run: func(_ *cobra.Command, _ []string) {
			if configuration.PKIBackend != nil {
				commands.VerifyPKIBackend(configuration)
			}
			// TODO: call install PKI monitor when implemented
		},
	}

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(verifyCmd)

	return rootCmd
}

func setUpGlobalFlags(cmd *cobra.Command, configFile *string) {
	flags := cmd.PersistentFlags()

	flags.StringVarP(
		configFile,
		"configFile",
		"f",
		"vvw_config.hcl",
		"Path to config file to use to configure Venafi Vault plugin",
	)
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
