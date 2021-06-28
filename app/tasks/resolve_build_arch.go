package tasks

import (
	"fmt"
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
)

func ResolveBuildArch(sshClient ssh.VaultSSHClient, pluginBuildArch string) error {
	var sshBuildArch string
	var definedBuildArch string

	if pluginBuildArch == "" {
		definedBuildArch = "linux"
	} else {
		definedBuildArch = pluginBuildArch
	}
	osType, arch, err := sshClient.CheckOSArch()
	if err != nil {
		return fmt.Errorf("unable to resolve client erch via SSH: %s", err)
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
		return fmt.Errorf("defined build architecture (%s) doesn't match client architecture (%s)", definedBuildArch, sshBuildArch)
	}

	return nil
}
