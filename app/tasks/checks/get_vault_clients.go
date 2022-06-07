package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/lib"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func GetAPIClient(section reporter.Section, apiAddress, token string) (api.VaultAPIClient, error) {
	check := section.AddCheck("Checking Vault API connection parameters...")

	vaultClient, err := api.NewClient(
		&api.Config{
			APIAddress: apiAddress,
			Token:      token,
		},
		lib.NewVaultAPI(),
	)
	if err != nil {
		check.Errorf("Error setting the Vault address for the Vault API client: %s", err)
		return nil, err
	}

	_, err = vaultClient.GetVaultConfig()
	if err != nil {
		check.Errorf("Error connecting to Vault API at %s and reading config: %s", apiAddress, err)
		return nil, err
	}

	check.Success("Connected to Vault API")
	return vaultClient, nil
}

func GetSSHClients(section reporter.Section, cfg *config.VaultConfig) ([]ssh.VaultSSHClient, func(), error) {
	check := section.AddCheck("Checking Vault SSH connection parameters...")

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
			check.Errorf("Error connecting to Vault server at %s over SSH: %s", s.Hostname, err)
			closeFunc()
			return nil, nil, err
		}
		closeFuncs = append(closeFuncs, func() {
			_ = sshClient.Close()
		})
		sshClients = append(sshClients, sshClient)
	}

	if len(sshClients) > 0 {
		check.Success("Connected to Vault via SSH")
	} else {
		check.Success("No SSH parameters given, skipped check")
	}

	return sshClients, closeFunc, nil
}
