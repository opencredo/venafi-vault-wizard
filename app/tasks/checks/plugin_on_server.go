package checks

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func InstallPluginOnServer(
	reportSection reporter.Section,
	sshClient ssh.VaultSSHClient,
	filepath string,
	pluginBytes []byte,
) error {
	check := reportSection.AddCheck("Copying plugin to Vault server...")
	err := sshClient.WriteFile(bytes.NewReader(pluginBytes), filepath)
	if err != nil {
		if errors.Is(err, ssh.ErrFileBusy) {
			check.Warning("File already exists and is busy, so cannot overwrite")
			return nil
		}

		check.Errorf("Error copying plugin to Vault: %s", err)
		return err
	}

	check.Success("Plugin copied to Vault server")

	return nil
}

func VerifyPluginOnServer(reportSection reporter.Section, sshClient ssh.VaultSSHClient, filepath string) error {
	pluginExistsCheck := reportSection.AddCheck("Checking plugin binary exists...")
	exists, err := sshClient.FileExists(filepath)
	if err != nil {
		pluginExistsCheck.Errorf("Error checking plugin binary exists: %s", err)
		return err
	}
	if !exists {
		pluginExistsCheck.Errorf("Plugin binary does not exist at %s on Vault server", filepath)
		return fmt.Errorf("plugin not found on server")
	}
	pluginExistsCheck.Success("Found plugin binary on Vault server")

	return nil
}
