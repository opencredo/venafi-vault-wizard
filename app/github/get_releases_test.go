package github

import "testing"

const (
	venafiPKIBackendRepo = "Venafi/vault-pki-backend-venafi"
	venafiPKIMonitorRepo = "Venafi/vault-pki-monitor-venafi"
)

func TestGetReleases(t *testing.T) {
	tests := map[string]struct {
		repoOwnerAndName   string
		desiredVersion     string
		assetSearchSubstr  string
		expectedReleaseURL string
		wantErr            bool
	}{
		"venafi-pki-backend v0.9.0": {
			repoOwnerAndName:   venafiPKIBackendRepo,
			desiredVersion:     "v0.9.0",
			assetSearchSubstr:  "linux.zip",
			expectedReleaseURL: "https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v0.9.0/venafi-pki-backend_v0.9.0_linux.zip",
			wantErr:            false,
		},
		"venafi-pki-backend wrong version": {
			repoOwnerAndName:   venafiPKIBackendRepo,
			desiredVersion:     "0.9.0",
			assetSearchSubstr:  "linux.zip",
			expectedReleaseURL: "",
			wantErr:            true,
		},
		"venafi-pki-monitor strict v0.9.0": {
			repoOwnerAndName:   venafiPKIMonitorRepo,
			desiredVersion:     "v0.9.0",
			assetSearchSubstr:  "linux_strict.zip",
			expectedReleaseURL: "https://github.com/Venafi/vault-pki-monitor-venafi/releases/download/v0.9.0/venafi-pki-monitor_v0.9.0_linux_strict.zip",
			wantErr:            false,
		},
		"venafi-pki-monitor wrong asset search string": {
			repoOwnerAndName:   venafiPKIMonitorRepo,
			desiredVersion:     "v0.9.0",
			assetSearchSubstr:  "linux_strictly.zip",
			expectedReleaseURL: "",
			wantErr:            true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			url, version, err := GetReleases(
				test.repoOwnerAndName,
				test.desiredVersion,
				test.assetSearchSubstr,
			)
			if (err != nil) != test.wantErr {
				t.Errorf("GetReleases() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			// don't bother checking url if we wanted an error
			if test.wantErr {
				return
			}

			if url != test.expectedReleaseURL {
				t.Errorf("GetReleases mismatch, want %s, got %s", test.expectedReleaseURL, url)
			}
			if version != test.desiredVersion {
				t.Errorf("GetReleases mismatch, want %s, got %s", test.desiredVersion, version)
			}
		})
	}
}
