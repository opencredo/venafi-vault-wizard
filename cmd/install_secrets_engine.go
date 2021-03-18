package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter/pretty"
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
	report := pretty.NewReport()

	vaultURL, err := url.Parse(vaultAddress)
	if err != nil {
		// TODO: report error better. maybe a reporter.Section or Check for getting these clients set up
		fmt.Println("Invalid Vault Address")
		return
	}

	vaultSSHAddress := fmt.Sprintf("%s:%d", vaultURL.Hostname(), sshPort)

	sshClient, err := ssh.NewClient(vaultSSHAddress, sshUser, sshPassword)
	if err != nil {
		// TODO: check errors here and report better than just a simple print
		fmt.Println("Error making SSH connection")
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
		Reporter:        report,
		PluginURL:       pluginURL,
		PluginName:      pluginName,
		PluginMountPath: pluginMountPath,
	})
	if err != nil {
		return
	}

	err = tasks.ConfigureVenafiPKIBackend(&tasks.ConfigureVenafiPKIBackendInput{
		VaultClient:     vaultClient,
		Reporter:        report,
		PluginMountPath: pluginMountPath,
		SecretName:      "cloud", // TODO: override on command line
		RoleName:        "cloud",
		VenafiAPIKey:    venafiAPIKey,
		VenafiZoneID:    venafiZoneID,
	})
	if err != nil {
		return
	}

	report.Finish(
		fmt.Sprintf(
			"Finished! You can try and request a certificate using:\n$ vault write %s/issue/%s common_name=\"test.example.com\"\n",
			pluginMountPath,
			"cloud",
		),
		"Success! Vault is configured to work with Venafi",
	)
}
