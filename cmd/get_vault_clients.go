package cmd

import (
	"fmt"
	"net/url"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/lib"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func getClients(report reporter.Report) (ssh.VaultSSHClient, api.VaultAPIClient, func(), error) {
	checkConnectionSection := report.AddSection("Checking connection to Vault")
	check := checkConnectionSection.AddCheck("Checking Vault connection parameters...")

	vaultURL, err := url.Parse(vaultAddress)
	if err != nil {
		check.Error(fmt.Sprintf("Invalid Vault address: %s", err))
		return nil, nil, nil, err
	}

	vaultClient := api.NewClient(
		&api.Config{
			APIAddress: vaultAddress,
			Token:      vaultToken,
		},
		lib.NewVaultAPI(),
	)
	_, err = vaultClient.GetVaultConfig()
	if err != nil {
		check.Error(fmt.Sprintf("Error connecting to Vault API and reading config: %s", err))
		return nil, nil, nil, err
	}

	check.UpdateStatus("Successfully connected to Vault API, establishing SSH connection...")

	vaultSSHAddress := fmt.Sprintf("%s:%d", vaultURL.Hostname(), sshPort)
	sshClient, err := ssh.NewClient(vaultSSHAddress, sshUser, sshPassword)
	if err != nil {
		check.Error(fmt.Sprintf("Error connecting to Vault server over SSH: %s", err))
		return nil, nil, nil, err
	}
	closeFunc := func() {
		_ = sshClient.Close()
	}

	check.Success("Connected to Vault via its API and SSH")

	return sshClient, vaultClient, closeFunc, nil
}
