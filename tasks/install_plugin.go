package tasks

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/pterm/pterm"

	"github.com/opencredo/venafi-vault-wizard/helpers/download_plugin"
	"github.com/opencredo/venafi-vault-wizard/helpers/vault"
)

type InstallPluginInput struct {
	VaultClient     vault.Vault
	Downloader      download_plugin.PluginDownloader
	PluginURL       string
	PluginName      string
	PluginMountPath string
}

func InstallPlugin(input *InstallPluginInput) error {
	pterm.DefaultSection.Println("Checking Vault")

	vaultCheckSpinner, _ := pterm.DefaultSpinner.Start("Checking Vault...")
	vaultCheckSpinner.UpdateText("Checking Vault Plugin Directory...")

	pluginDir, err := input.VaultClient.GetPluginDir()
	if err != nil {
		if errors.Is(err, vault.ErrReadingVaultPath) {
			vaultCheckSpinner.Fail(fmt.Sprintf("Error getting plugin directory: %s", err))
		} else if errors.Is(err, vault.ErrPluginDirNotConfigured) {
			vaultCheckSpinner.Fail("The plugin_directory hasn't been configured correctly in the Vault Server Config")
		}
		return err
	}
	vaultCheckSpinner.Success("Vault Plugin Directory found")

	pterm.Println()
	pterm.Printf("The Vault server plugin directory is configured as %s\n", pluginDir)

	pterm.DefaultSection.Println("Installing plugin to Vault")

	pluginPath := fmt.Sprintf("%s/%s", pluginDir, input.PluginName)

	sha, err := installPlugin(input.VaultClient, input.Downloader, pluginPath, input.PluginURL)
	if err != nil {
		return err
	}

	err = enablePluginMlock(input.VaultClient, pluginPath)
	if err != nil {
		return err
	}

	err = enablePlugin(input.VaultClient, input.PluginName, sha)
	if err != nil {
		return err
	}

	pterm.Println()
	pterm.Printf("The Venafi plugin has been installed as %s\n", input.PluginName)

	pterm.DefaultSection.Println("Configuring plugin")

	err = mountPlugin(input.VaultClient, input.PluginName, input.PluginMountPath)
	if err != nil {
		return err
	}

	pterm.Println()
	pterm.Printf("The plugin has been mounted at as %s\n", input.PluginMountPath)

	return nil
}

func installPlugin(client vault.Vault, downloader download_plugin.PluginDownloader, pluginPath, pluginURL string) (string, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Installing plugin to Vault server...")
	plugin, sha, err := downloader.DownloadPluginAndUnzip(pluginURL)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
		return "", err
	}

	spinner.UpdateText("Successfully downloaded plugin, copying to Vault server over SSH...")

	err = client.WriteFile(bytes.NewReader(plugin), pluginPath)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error copying plugin to Vault: %s", err))
		return "", err
	}

	spinner.Success("Plugin copied to Vault server")

	return sha, nil
}

func enablePluginMlock(client vault.Vault, pluginFilepath string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Checking if mlock is disabled...")
	mlockDisabled, err := client.IsMLockDisabled()
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error checking whether mlock is disabled: %s", err))
		return err
	}

	if !mlockDisabled {
		spinner.UpdateText("Mlock is enabled on the Vault server, attempting to add IPC_LOCK capability to plugin...")

		err := client.AddIPCLockCapabilityToFile(pluginFilepath)
		if err != nil {
			spinner.Warning("Error adding IPC_LOCK capability to plugin, might be needed for mlock: %s", err)
			return nil
		}

		spinner.Success("IPC_LOCK capability added to plugin")
		return nil
	}

	spinner.Warning("Mlock is disabled on the Vault server, should be enabled for production")
	return nil
}

func enablePlugin(client vault.Vault, pluginName, sha string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Enabling plugin in Vault plugin catalog...")
	err := client.RegisterPlugin(pluginName, pluginName, sha)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error registering plugin in Vault catalog: %s", err))
		return err
	}

	spinner.Success("Successfully registered plugin in Vault plugin catalog")
	return nil
}

func mountPlugin(client vault.Vault, pluginName, mountPath string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Mounting plugin ...")
	err := client.MountPlugin(pluginName, mountPath)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error mounting plugin: %s", err))
		return err
	}

	spinner.Success("Plugin mounted")
	return nil
}
