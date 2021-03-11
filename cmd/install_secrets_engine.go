package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/opencredo/venafi-vault-wizard/helpers"
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

const pluginURL = "https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v0.8.3/venafi-pki-backend_v0.8.3_linux.zip"

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
	vaultCheckSpinner.Success("Vault Plugin Directory found")

	pterm.Println()
	pterm.Printf("The Vault server plugin directory is correctly configured as %s\n", pluginDir)

	pterm.DefaultSection.Println("Installing plugin to Vault")

	pluginInstallSpinner, _ := pterm.DefaultSpinner.Start("Installing plugin to Vault server...")
	plugin, sha, err := helpers.DownloadPluginAndUnzip(pluginURL)
	if err != nil {
		pluginInstallSpinner.Fail(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
		return
	}

	pluginInstallSpinner.UpdateText("Successfully downloaded plugin, copying to Vault server over SSH...")

	err = vaultClient.WriteFile(bytes.NewReader(plugin), fmt.Sprintf("%s/venafi-pki-backend", pluginDir))
	if err != nil {
		pluginInstallSpinner.Fail(fmt.Sprintf("Error copying plugin to Vault: %s", err))
		return
	}
	pluginInstallSpinner.Success("Plugin copied to Vault server")

	const pluginName = "venafi-pki-backend"

	pluginEnableSpinner, _ := pterm.DefaultSpinner.Start("Enabling plugin in Vault plugin catalog...")
	err = vaultClient.RegisterPlugin(pluginName, pluginName, sha)
	if err != nil {
		pluginEnableSpinner.Fail(fmt.Sprintf("Error registering plugin in Vault catalog: %s", err))
		return
	}
	pluginEnableSpinner.Success("Successfully registered plugin in Vault plugin catalog")

	pterm.Println()
	pterm.Printf("The Venafi plugin has been installed as %s\n", pluginName)

	pterm.Println()
	pterm.DefaultHeader.Println("Success! Vault is configured to work with Venafi")
}
