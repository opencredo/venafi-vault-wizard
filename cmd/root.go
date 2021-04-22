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

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies desired state as specified in config file",
		Long:  "Reads the config file and makes necessary changes to Vault server(s) specified to install and configure plugin(s) specified",
		Run: func(_ *cobra.Command, _ []string) {
			commands.Apply(configuration)
		},
	}

	rootCmd.AddCommand(applyCmd)

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
