package tasks

import (
	"github.com/opencredo/venafi-vault-wizard/app/vault/ssh"
	mockReport "github.com/opencredo/venafi-vault-wizard/mocks/app/reporter"
	mockSSH "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/ssh"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestResolveBuildArch(t *testing.T) {
	vaultSSHClient := new(mockSSH.VaultSSHClient)
	report := new(mockReport.Report)
	section := new(mockReport.Section)
	check := new(mockReport.Check)
	defer vaultSSHClient.AssertExpectations(t)
	defer report.AssertExpectations(t)
	defer section.AssertExpectations(t)
	defer check.AssertExpectations(t)

	reportExpectations(report, section, check)

	tests := map[string]struct {
		pluginBuildArch string
		sshOSType       string
		sshArch         string
		wantErr         bool
	}{
		"valid darwin 64bit": {
			pluginBuildArch: "darwin",
			sshOSType:       "Darwin",
			sshArch:         "x86_64",
			wantErr:         false,
		},
		"valid linux 64bit": {
			pluginBuildArch: "linux",
			sshOSType:       "Linux",
			sshArch:         "x86_64",
			wantErr:         false,
		},
		"valid linux 32bit": {
			pluginBuildArch: "linux86",
			sshOSType:       "Linux",
			sshArch:         "i686",
			wantErr:         false,
		},
		"invalid linux 64bit": {
			pluginBuildArch: "linux86",
			sshOSType:       "Linux",
			sshArch:         "x86_64",
			wantErr:         true,
		},
	}

	check.On("Errorf", mock.AnythingOfType("string"), mock.Anything).Return().Maybe()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			vaultSSHClient.On("CheckOSArch").Return(test.sshOSType, test.sshArch, nil).Once()
			err := ResolveBuildArch(&ResolveBuildArchInput{
				SSHClients:      []ssh.VaultSSHClient{vaultSSHClient},
				PluginBuildArch: test.pluginBuildArch,
				Reporter:        report,
			})

			if (err != nil) != test.wantErr {
				t.Errorf("expected no error, got: '%s'", err)
				return
			}
		})
	}
}
