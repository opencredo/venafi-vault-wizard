package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type ResolveBuildArchInput struct {
	SSHClients      []ssh.VaultSSHClient
	PluginBuildArch string
	Reporter        reporter.Report
}

func ResolveBuildArch(input *ResolveBuildArchInput) error {
	buildArchSection := input.Reporter.AddSection("Checking Vault Server OS and CPU Architecture")
	var sshBuildArch string
	var definedBuildArch string

	if input.PluginBuildArch == "" {
		definedBuildArch = "linux"
	} else {
		definedBuildArch = input.PluginBuildArch
	}

	for i, sshClient := range input.SSHClients {
		check := buildArchSection.AddCheck(fmt.Sprintf("Checking Vault Server %d", i+1))
		osType, arch, err := sshClient.CheckOSArch()
		if err != nil {
			check.Errorf("Unable to resolve client arch via SSH: %w", err)
			return err
		}

		switch osType {
		case "Darwin":
			sshBuildArch = "darwin"
		case "Linux":
			sshBuildArch = "linux"
		}

		if arch != "x86_64" {
			sshBuildArch = sshBuildArch + "86"
		}

		if definedBuildArch != sshBuildArch {
			check.Errorf("Defined build architecture (%s) doesn't match client architecture (%s)", definedBuildArch, sshBuildArch)
			return fmt.Errorf("Defined build architecture (%s) doesn't match client architecture (%s)", definedBuildArch, sshBuildArch)
		}

		check.Successf("Requested plugin build architecture (%s) matches Vault Server %d", definedBuildArch, i+1)
	}

	return nil
}
