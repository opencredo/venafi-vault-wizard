package checks

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func InstallPluginOnServer(
	reportSection reporter.Section,
	sshClient ssh.VaultSSHClient,
	download downloader.PluginDownloader,
	filepath string,
	pluginURL string,
) (string, error) {
	check := reportSection.AddCheck("Installing plugin to Vault server...")
	plugin, sha, err := download.DownloadPluginAndUnzip(pluginURL)
	if err != nil {
		check.Error(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
		return "", err
	}

	check.UpdateStatus("Successfully downloaded plugin, copying to Vault server over SSH...")

	err = sshClient.WriteFile(bytes.NewReader(plugin), filepath)
	if err != nil {
		if errors.Is(err, ssh.ErrFileBusy) {
			check.Warning("File already exists and is busy, so cannot overwrite")
			return sha, nil
		}

		check.Error(fmt.Sprintf("Error copying plugin to Vault: %s", err))
		return "", err
	}

	check.Success("Plugin copied to Vault server")

	return sha, nil
}

func VerifyPluginOnServer(reportSection reporter.Section, sshClient ssh.VaultSSHClient, filepath string) error {
	pluginExistsCheck := reportSection.AddCheck("Checking plugin binary exists...")
	exists, err := sshClient.FileExists(filepath)
	if err != nil {
		pluginExistsCheck.Error(fmt.Sprintf("Error checking plugin binary exists: %s", err))
		return err
	}
	if !exists {
		pluginExistsCheck.Error(fmt.Sprintf("Plugin binary does not exist at %s on Vault server", filepath))
		return fmt.Errorf("plugin not found on server")
	}
	pluginExistsCheck.Success("Found plugin binary on Vault server")

	return nil
}
