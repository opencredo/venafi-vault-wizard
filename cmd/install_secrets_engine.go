package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/tasks"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/lib"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
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

	// TODO: get from command-line
	sshClient, err := ssh.NewClient(sshAddress, "vagrant", "vagrant")
	if err != nil {
		return
	}
	defer sshClient.Close()

	vaultClient := api.NewClient(
		&api.Config{
			APIAddress: vaultAddress,
			Token:      vaultToken,
		},
		lib.NewVaultAPI(),
	)

	err = tasks.InstallPlugin(&tasks.InstallPluginInput{
		VaultClient:     vaultClient,
		SSHClient:       sshClient,
		Downloader:      downloader.NewPluginDownloader(),
		PluginURL:       pluginURL,
		PluginName:      pluginName,
		PluginMountPath: pluginMountPath,
	})
	if err != nil {
		return
	}

	err = tasks.ConfigureVenafiPKIBackend(&tasks.ConfigureVenafiPKIBackendInput{
		VaultClient:     vaultClient,
		PluginMountPath: pluginMountPath,
		SecretName:      "cloud", // TODO: override on command line
		RoleName:        "cloud",
		VenafiAPIKey:    venafiAPIKey,
		VenafiZoneID:    venafiZoneID,
	})
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
