package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/opencredo/venafi-vault-wizard/app/questions/prompter"
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/commands"
)

func NewRootCommand() *cobra.Command {
	var configFile string

	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "vvw",
		Short: "Venafi Vault Wizard",
		Long:  "VVW is a wizard to automate the installation and verification of Venafi PKI plugins for HashiCorp Vault.",
	}

	setUpGlobalFlags(rootCmd, &configFile)

	generateConfigCmd := &cobra.Command{
		Use:   "generate-config",
		Short: "Generates config file based on asking questions",
		Long:  "Builds up a config file suitable for use with the apply command by asking questions about the Vault setup",
		Run: func(_ *cobra.Command, _ []string) {
			questioner := prompter.NewPrompter()
			commands.GenerateConfig(configFile, questioner)
		},
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies desired state as specified in config file",
		Long:  "Reads the config file and makes necessary changes to Vault server(s) specified to install and configure plugin(s) specified",
		RunE: func(_ *cobra.Command, _ []string) error {
			src, err := ioutil.ReadFile(configFile)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("error: config file %s not found: %s", configFile, err)
				}

				return fmt.Errorf("can't read %s: %w", configFile, err)
			}

			return commands.Apply(configFile, src)
		},
	}

	rootCmd.AddCommand(generateConfigCmd)
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
