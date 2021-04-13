package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewConfig(t *testing.T) {
	tests := map[string]struct {
		config  string
		want    *Config
		wantErr bool
	}{
		"valid venafi-pki-backend with Cloud": {
			config:  validPKIBackendCloudConfig,
			want:    validPKIBackendCloudConfigResult,
			wantErr: false,
		},
		"valid venafi-pki-backend with TPP": {
			config:  validPKIBackendTPPConfig,
			want:    validPKIBackendTPPConfigResult,
			wantErr: false,
		},
		"invalid venafi-pki-backend without role": {
			config:  invalidPKIBackendConfigNoRole,
			wantErr: true,
		},
		"invalid venafi-pki-backend without secret": {
			config:  invalidPKIBackendConfigNoSecret,
			wantErr: true,
		},
		"invalid venafi-pki-backend with both cloud and TPP": {
			config:  invalidPKIBackendConfigBothConnectionTypes,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewConfig("vvwconfig.hcl", []byte(tt.config))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("NewConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

const validPKIBackendCloudConfig = `
vault {
  address = "http://localhost:8200"
  token = "root"

  ssh {
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

venafi_pki_backend {
  mount_path = "venafi-pki"

  role "cloud" {
    secret "cloud" {
      venafi_cloud {
        apikey = "apikey"
        zone = "zone"
      }
    }

    test_certificate {
      common_name = "vvw-example.test"
    }
  }
}
`

var validPKIBackendCloudConfigResult = &Config{
	Vault: VaultConfig{
		VaultAddress: "http://localhost:8200",
		VaultToken:   "root",
		SSHConfig: SSH{
			Username: "vagrant",
			Password: "vagrant",
			Port:     22,
		},
	},
	PKIBackend: &VenafiPKIBackendConfig{
		MountPath: "venafi-pki",
		Roles: []Role{
			{
				Name: "cloud",
				Secret: VenafiSecret{
					Name: "cloud",
					Cloud: &VenafiCloudConnection{
						APIKey: "apikey",
						Zone:   "zone",
					},
				},
				TestCerts: []CertificateRequest{
					{
						CommonName: "vvw-example.test",
					},
				},
			},
		},
	},
}

const validPKIBackendTPPConfig = `
vault {
  address = "http://localhost:8200"
  token = "root"

  ssh {
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

venafi_pki_backend {
  mount_path = "venafi-pki"

  role "tppRole" {
    secret "tpptest" {
      venafi_tpp {
        url = "tpp.venafitest.com"
        username = "admin"
        password = "pword234"
        policy = "Partner Dev\\\\TLS\\\\HashiCorp Vault"
      }
    }

    test_certificate {
      common_name = "vvw-example.test"
    }
  }
}`

var validPKIBackendTPPConfigResult = &Config{
	Vault: VaultConfig{
		VaultAddress: "http://localhost:8200",
		VaultToken:   "root",
		SSHConfig: SSH{
			Username: "vagrant",
			Password: "vagrant",
			Port:     22,
		},
	},
	PKIBackend: &VenafiPKIBackendConfig{
		MountPath: "venafi-pki",
		Roles: []Role{
			{
				Name: "tppRole",
				Secret: VenafiSecret{
					Name: "tpptest",
					TPP: &VenafiTPPConnection{
						URL:      "tpp.venafitest.com",
						Username: "admin",
						Password: "pword234",
						Policy:   "Partner Dev\\\\TLS\\\\HashiCorp Vault",
					},
				},
				TestCerts: []CertificateRequest{
					{
						CommonName: "vvw-example.test",
					},
				},
			},
		},
	},
}

const invalidPKIBackendConfigNoRole = `
vault {
  address = "http://localhost:8200"
  token = "root"

  ssh {
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

venafi_pki_backend {
  mount_path = "venafi-pki"
}`

const invalidPKIBackendConfigNoSecret = `
vault {
  address = "http://localhost:8200"
  token = "root"

  ssh {
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

venafi_pki_backend {
  mount_path = "venafi-pki"

  role "cloud" {
  }
}`

const invalidPKIBackendConfigBothConnectionTypes = `
vault {
  address = "http://localhost:8200"
  token = "root"

  ssh {
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

venafi_pki_backend {
  mount_path = "venafi-pki"

  role "cloud" {
    secret "cloud" {
      venafi_cloud {
        apikey = "apikey"
        zone = "zone"
      }
      venafi_tpp {
        url = "tpp.venafitest.com"
        username = "admin"
        password = "pword234"
        policy = "Partner Dev\\TLS\\HashiCorp Vault"
      }
    }
  }
}`
