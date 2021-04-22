package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/lib"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func GetClients(cfg *config.VaultConfig, report reporter.Report) ([]ssh.VaultSSHClient, api.VaultAPIClient, func(), error) {
	checkConnectionSection := report.AddSection("Checking connection to Vault")
	check := checkConnectionSection.AddCheck("Checking Vault connection parameters...")

	vaultClient := api.NewClient(
		&api.Config{
			APIAddress: cfg.VaultAddress,
			Token:      cfg.VaultToken,
		},
		lib.NewVaultAPI(),
	)
	_, err := vaultClient.GetVaultConfig()
	if err != nil {
		check.Error(fmt.Sprintf("Error connecting to Vault API at %s and reading config: %s", cfg.VaultAddress, err))
		return nil, nil, nil, err
	}

	check.UpdateStatus("Successfully connected to Vault API, establishing SSH connection...")

	var sshClients []ssh.VaultSSHClient
	var closeFuncs []func()
	closeFunc := func() {
		for _, f := range closeFuncs {
			f()
		}
	}

	for _, s := range cfg.SSHConfig {
		address := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
		sshClient, err := ssh.NewClient(address, s.Username, s.Password)
		if err != nil {
			check.Error(fmt.Sprintf("Error connecting to Vault server at %s over SSH: %s", s.Hostname, err))
			closeFunc()
			return nil, nil, nil, err
		}
		closeFuncs = append(closeFuncs, func() {
			_ = sshClient.Close()
		})
		sshClients = append(sshClients, sshClient)
	}

	check.Success("Connected to Vault via its API and SSH")

	return sshClients, vaultClient, closeFunc, nil
}
