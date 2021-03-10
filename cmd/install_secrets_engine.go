package cmd

import (
	"errors"
	"fmt"
	"github.com/opencredo/venafi-vault-wizard/vault"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var installPKIBackendCommand = &cobra.Command{
	Use:   "venafi-pki-backend",
	Short: "Installs venafi-pki-backend plugin",
	Long:  "Installs the venafi-pki-backend plugin to allow Vault to request certificates from Venafi",
	Run:   installPKIBackend,
}

func init() {
	installCommand.AddCommand(installPKIBackendCommand)
}

func installPKIBackend(_ *cobra.Command, _ []string) {
	pterm.Error.ShowLineNumber = false

	pterm.DefaultSection.Println("Checking Vault")

	vaultCheckSpinner, _ := pterm.DefaultSpinner.Start("Checking Vault...")
	vaultClient, err := vault.NewVault(&vault.Config{
		APIAddress: vaultAddress,
		Token:      vaultToken,
		SSHAddress: sshAddress,
	})
	if err != nil {
		vaultCheckSpinner.Fail(fmt.Sprintf("Error with Vault parameters: %s", err))
		return
	}

	vaultCheckSpinner.UpdateText("Checking Vault Plugin Directory...")

	pluginDir, err := vaultClient.GetPluginDir()
	if err != nil {
		if errors.Is(err, vault.ErrReadingVaultPath) {
			vaultCheckSpinner.Fail(fmt.Sprintf("Error getting plugin directory: %s", err))
		} else if errors.Is(err, vault.ErrPluginDirNotConfigured) {
			vaultCheckSpinner.Fail("The plugin_directory hasn't been configured correctly in the Vault Server Config")
		}
		return
	}
	vaultCheckSpinner.Success("Vault Plugin Directory found!")

	pterm.Println()
	pterm.Printf("The Vault server plugin directory is correctly configured as %s\n", pluginDir)

	pterm.DefaultSection.Println("Installing plugin to Vault")

	pluginInstallSpinner, _ := pterm.DefaultSpinner.Start("Installing plugin to Vault server...")
	err = vaultClient.CopyFile("/config/randomfile.txt")
	if err != nil {
		pluginInstallSpinner.Fail(fmt.Sprintf("Error copying file %s", err))
		return
	}
	pluginInstallSpinner.Success("Plugin Installed!")

	pterm.DefaultHeader.Println("Success! Vault is configured to work with Venafi")
}
