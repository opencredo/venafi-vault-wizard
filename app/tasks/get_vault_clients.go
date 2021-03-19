package tasks

import (
	"fmt"
	"net/url"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/lib"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func GetClients(cfg *config.GlobalConfig, report reporter.Report) (ssh.VaultSSHClient, api.VaultAPIClient, func(), error) {
	checkConnectionSection := report.AddSection("Checking connection to Vault")
	check := checkConnectionSection.AddCheck("Checking Vault connection parameters...")

	vaultURL, err := url.Parse(cfg.VaultAddress)
	if err != nil {
		check.Error(fmt.Sprintf("Invalid Vault address: %s", err))
		return nil, nil, nil, err
	}

	vaultClient := api.NewClient(
		&api.Config{
			APIAddress: cfg.VaultAddress,
			Token:      cfg.VaultToken,
		},
		lib.NewVaultAPI(),
	)
	_, err = vaultClient.GetVaultConfig()
	if err != nil {
		check.Error(fmt.Sprintf("Error connecting to Vault API at %s and reading config: %s", cfg.VaultAddress, err))
		return nil, nil, nil, err
	}

	check.UpdateStatus("Successfully connected to Vault API, establishing SSH connection...")

	vaultSSHAddress := fmt.Sprintf("%s:%d", vaultURL.Hostname(), cfg.SSHPort)
	sshClient, err := ssh.NewClient(vaultSSHAddress, cfg.SSHUser, cfg.SSHPassword)
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
