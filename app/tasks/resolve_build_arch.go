package tasks

import (
	"fmt"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

type ResolveBuildArchInput struct {
	SSHClients []ssh.VaultSSHClient
	PluginBuildArch string
	Reporter reporter.Report
}

func ResolveBuildArch(input *ResolveBuildArchInput) error {
	var sshBuildArch string
	var definedBuildArch string

	if input.PluginBuildArch == "" {
		definedBuildArch = "linux"
	} else {
		definedBuildArch = input.PluginBuildArch
	}

	for i, sshClient := range input.SSHClients {
		checkSSHClient := input.Reporter.AddSection(fmt.Sprintf("Checking SSH Client #%d", i))

		check := checkSSHClient.AddCheck("Checking OS Architecture")
		osType, arch, err := sshClient.CheckOSArch()
		if err != nil {
			check.Errorf("unable to resolve client arch via SSH: %s", err)
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
			check.Errorf("defined build architecture (%s) doesn't match client architecture (%s)", definedBuildArch, sshBuildArch)
		}
	}

	return nil
}
