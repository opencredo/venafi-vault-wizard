package cmd

import (
	"fmt"
	"github.com/opencredo/venafi-vault-wizard/tasks"
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

	vaultClient, err := vault.NewVault(&vault.Config{
		APIAddress: vaultAddress,
		Token:      vaultToken,
		SSHAddress: sshAddress,
	})
	if err != nil {
		return
	}

	err = tasks.InstallPlugin(&tasks.InstallPluginInput{
		VaultClient:     vaultClient,
		PluginURL:       pluginURL,
		PluginName:      pluginName,
		PluginMountPath: pluginMountPath,
	})
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
