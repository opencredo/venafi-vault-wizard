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

const (
	pluginName      = "venafi-pki-backend"
	pluginURL       = "https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v0.8.3/venafi-pki-backend_v0.8.3_linux.zip"
	pluginMountPath = "venafi-pki"
)

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

	sha, err := installPlugin(vaultClient, pluginDir)
	if err != nil {
		return
	}

	err = enablePlugin(vaultClient, sha)
	if err != nil {
		return
	}

	pterm.Println()
	pterm.Printf("The Venafi plugin has been installed as %s\n", pluginName)

	pterm.DefaultSection.Println("Configuring plugin")

	err = mountPlugin(vaultClient, pluginMountPath)
	if err != nil {
		return
	}

	pterm.Println()
	pterm.Printf("The plugin has been mounted at as %s\n", pluginMountPath)

	err = addVenafiSecret(vaultClient, "cloud")
	if err != nil {
		return
	}

	err = addVenafiRole(vaultClient, "cloud", "cloud")
	if err != nil {
		return
	}

	pterm.Println()
	pterm.DefaultBasicText.WithStyle(&pterm.Style{pterm.FgGreen}).
		Printf(
			"Finished! You can try and request a certificate using:\n"+
				"$ vault write %s/issue/%s common_name=\"test.example.com\"\n", pluginMountPath, "cloud")

	pterm.Println()
	pterm.DefaultHeader.Println("Success! Vault is configured to work with Venafi")
}

func installPlugin(client vault.Vault, pluginDir string) (string, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Installing plugin to Vault server...")
	plugin, sha, err := helpers.DownloadPluginAndUnzip(pluginURL)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
		return "", err
	}

	spinner.UpdateText("Successfully downloaded plugin, copying to Vault server over SSH...")

	err = client.WriteFile(bytes.NewReader(plugin), fmt.Sprintf("%s/venafi-pki-backend", pluginDir))
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error copying plugin to Vault: %s", err))
		return "", err
	}

	spinner.Success("Plugin copied to Vault server")
	return sha, nil
}

func enablePlugin(client vault.Vault, sha string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Enabling plugin in Vault plugin catalog...")
	err := client.RegisterPlugin(pluginName, pluginName, sha)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error registering plugin in Vault catalog: %s", err))
		return err
	}

	spinner.Success("Successfully registered plugin in Vault plugin catalog")
	return nil
}

func mountPlugin(client vault.Vault, mountPath string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Mounting plugin ...")
	err := client.MountPlugin(pluginName, mountPath)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error mounting plugin: %s", err))
		return err
	}

	spinner.Success("Plugin mounted")
	return nil
}

func addVenafiSecret(client vault.Vault, secretName string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Adding Venafi secret...")
	err := client.WriteValue("venafi-pki/venafi/"+secretName, map[string]interface{}{
		"apikey": venafiAPIKey,
		"zone":   venafiZoneID,
	})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error configuring Venafi secret: %s", err))
		return err
	}

	spinner.Success("Venafi secret configured at venafi-pki/venafi/" + secretName)
	return nil
}

func addVenafiRole(client vault.Vault, roleName, secretName string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Adding Venafi secret...")
	err := client.WriteValue("venafi-pki/roles/"+roleName, map[string]interface{}{
		"venafi_secret": secretName,
	})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error configuring Venafi role: %s", err))
		return err
	}

	spinner.Success("Venafi role configured")
	return nil
}
