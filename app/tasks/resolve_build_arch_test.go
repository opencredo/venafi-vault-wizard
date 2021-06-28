package tasks

import (
	"fmt"
	mockSSH "github.com/opencredo/venafi-vault-wizard/mocks/app/vault/ssh"
	"testing"
)

func TestResolveBuildArch(t *testing.T) {
	vaultSSHClient := new(mockSSH.VaultSSHClient)
	//report := new(mockReport.Report)
	//section := new(mockReport.Section)
	//check := new(mockReport.Check)
	defer vaultSSHClient.AssertExpectations(t)
	//defer report.AssertExpectations(t)
	//defer section.AssertExpectations(t)
	//defer check.AssertExpectations(t)

	//reportExpectations(report, section, check)

	tests := map[string]struct{
		pluginBuildArch string
		sshOSType string
		sshArch string
		expectedErr error
	}{
		"valid darwin 64bit": {
			pluginBuildArch: "darwin",
			sshOSType: "Darwin",
			sshArch: "x86_64",
			expectedErr: nil,
		},
		"valid linux 64bit": {
			pluginBuildArch: "linux",
			sshOSType: "Linux",
			sshArch: "x86_64",
			expectedErr: nil,
		},
		"valid linux 32bit": {
			pluginBuildArch: "linux86",
			sshOSType: "Linux",
			sshArch: "i686",
			expectedErr: nil,
		},
		"invalid linux 64bit": {
			pluginBuildArch: "linux86",
			sshOSType: "Linux",
			sshArch: "x86_64",
			expectedErr: fmt.Errorf("defined build architecture (linux86) doesn't match client architecture (linux)"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			vaultSSHClient.On("CheckOSArch").Return(test.sshOSType, test.sshArch, nil).Once()
			err := ResolveBuildArch(vaultSSHClient, test.pluginBuildArch)

			if test.expectedErr == nil && err != nil {
				t.Errorf("expected no error, got: '%s'", err)
				return
			}

			if test.expectedErr != nil && err != nil {
				if test.expectedErr.Error() != err.Error() {
					t.Errorf("expected '%s', got: '%s'", test.expectedErr, err)
					return
				} else {
					return
				}
			}
		})
	}
}

