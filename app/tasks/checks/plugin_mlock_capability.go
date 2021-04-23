package checks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func InstallPluginMlock(
	reportSection reporter.Section,
	sshClient ssh.VaultSSHClient,
	filepath string,
) error {
	check := reportSection.AddCheck("Attempting to add IPC_LOCK capability to plugin executable...")

	err := sshClient.AddIPCLockCapabilityToFile(filepath)
	if err != nil {
		check.Warning(fmt.Sprintf("Error adding IPC_LOCK capability to plugin, might be needed for mlock: %s", err))
		return nil
	}

	check.Success("IPC_LOCK capability added to plugin executable")
	return nil
}

func VerifyPluginMlock(
	reportSection reporter.Section,
	vaultClient api.VaultAPIClient,
	sshClient ssh.VaultSSHClient,
	filepath string,
) error {
	check := reportSection.AddCheck("Checking whether plugin needs the IPC_LOCK capability and has it...")
	mLockDisabled, err := vaultClient.IsMLockDisabled()
	if err != nil {
		check.Error(fmt.Sprintf("Error checking whether mlock is disabled: %s", err))
		return err
	}
	if mLockDisabled {
		check.Warning("Mlock is disabled on the Vault server, should be enabled for production")
	} else {
		capOnFile, err := sshClient.IsIPCLockCapabilityOnFile(filepath)
		if err != nil {
			check.Error(fmt.Sprintf("Error checking plugin binary for the IPC_LOCK capabiltiy: %s", err))
			return err
		}

		if capOnFile {
			check.Success("Mlock is enabled and the plugin binary has the IPC_LOCK capability")
		} else {
			check.Warning("Mlock is enabled on Vault server but the plugin does not appear to have the IPC_LOCK capability, however if things seem to work then ignore this warning")
		}
	}

	return nil
}
