package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
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

			// don't bother checking config if we wanted an error
			if tt.wantErr {
				return
			}

			deletePluginsUncheckedFields(got)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("NewConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

const validPKIBackendCloudConfig = `
vault {
  api_address = "http://localhost:8200"
  token = "root"

  ssh {
    hostname = "localhost"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-backend" "venafi-pki" {
  version = "v0.9.0"
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
		SSHConfig: []SSH{
			{
				Hostname: "localhost",
				Username: "vagrant",
				Password: "vagrant",
				Port:     22,
			},
		},
	},
	Plugins: []plugins.Plugin{
		{
			Type:      "venafi-pki-backend",
			MountPath: "venafi-pki",
			Version:   "v0.9.0",
			Config:    nil,
			Impl: &pki_backend.VenafiPKIBackendConfig{
				MountPath: "venafi-pki",
				Version:   "v0.9.0",
				Roles: []pki_backend.Role{
					{
						Name: "cloud",
						Secret: venafi.VenafiSecret{
							Name: "cloud",
							Cloud: &venafi.VenafiCloudConnection{
								APIKey: "apikey",
								Zone:   "zone",
							},
						},
						TestCerts: []pki_backend.CertificateRequest{
							{
								CommonName: "vvw-example.test",
							},
						},
					},
				},
			},
		},
	},
}

const validPKIBackendTPPConfig = `
vault {
  api_address = "http://localhost:8200"
  token = "root"

  ssh {
    hostname = "localhost"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-backend" "venafi-pki" {
  version = "v0.9.0"
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
		SSHConfig: []SSH{
			{
				Hostname: "localhost",
				Username: "vagrant",
				Password: "vagrant",
				Port:     22,
			},
		},
	},
	Plugins: []plugins.Plugin{
		{
			Type:      "venafi-pki-backend",
			MountPath: "venafi-pki",
			Version:   "v0.9.0",
			Config:    nil,
			Impl: &pki_backend.VenafiPKIBackendConfig{
				MountPath: "venafi-pki",
				Version:   "v0.9.0",
				Roles: []pki_backend.Role{
					{
						Name: "tppRole",
						Secret: venafi.VenafiSecret{
							Name: "tpptest",
							TPP: &venafi.VenafiTPPConnection{
								URL:      "tpp.venafitest.com",
								Username: "admin",
								Password: "pword234",
								Policy:   "Partner Dev\\\\TLS\\\\HashiCorp Vault",
							},
						},
						TestCerts: []pki_backend.CertificateRequest{
							{
								CommonName: "vvw-example.test",
							},
						},
					},
				},
			},
		},
	},
}

const invalidPKIBackendConfigNoRole = `
vault {
  api_address = "http://localhost:8200"
  token = "root"

  ssh {
    hostname = "localhost"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-backend" "venafi-pki" {
}`

const invalidPKIBackendConfigNoSecret = `
vault {
  api_address = "http://localhost:8200"
  token = "root"

  ssh {
    hostname = "localhost"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-backend" "venafi-pki" {
  version = "v0.9.0"
  role "cloud" {
  }
}`

const invalidPKIBackendConfigBothConnectionTypes = `
vault {
  api_address = "http://localhost:8200"
  token = "root"

  ssh {
    hostname = "localhost"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-backend" "venafi-pki" {
  version = "v0.9.0"

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

func deletePluginsUncheckedFields(config *Config) {
	for i := 0; i < len(config.Plugins); i++ {
		config.Plugins[i].Config = nil
	}
}
