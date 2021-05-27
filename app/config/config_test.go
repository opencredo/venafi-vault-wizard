package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
	pki_monitor "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-monitor"
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
		"valid venafi-pki-monitor with Cloud": {
			config:  validPKIMonitorConfig,
			want:    validPKIMonitorConfigResult,
			wantErr: false,
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
    zone = "zone"
    secret "cloud" {
      venafi_cloud {
        apikey = "apikey"
		zone = "zone1"
      }
    }

    test_certificate {
	  common_name = "vvw-example.test"
	  ou = "VVW"
	  organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "1h"
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
						Zone: "zone",
						Secret: venafi.VenafiSecret{
							Name: "cloud",
							Cloud: &venafi.VenafiCloudConnection{
								APIKey: "apikey",
								Zone:   "zone1",
							},
						},
						TestCerts: []venafi.CertificateRequest{
							{
								CommonName:   "vvw-example.test",
								OU:           "VVW",
								Organisation: "VVW",
								Locality:     "London",
								Province:     "London",
								Country:      "GB",
								TTL:          "1h",
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
    zone = "Partner Dev\\\\TLS\\\\HashiCorp Vault"
    secret "tpptest" {
      venafi_tpp {
        url = "tpp.venafitest.com"
        username = "admin"
        password = "pword234"
		zone = "zone1"
      }
    }

    test_certificate {
	  common_name = "vvw-example.test"
	  ou = "VVW"
	  organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "1h"
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
						Zone: "Partner Dev\\\\TLS\\\\HashiCorp Vault",
						Secret: venafi.VenafiSecret{
							Name: "tpptest",
							TPP: &venafi.VenafiTPPConnection{
								URL:      "tpp.venafitest.com",
								Username: "admin",
								Password: "pword234",
								Zone:     "zone1",
							},
						},
						TestCerts: []venafi.CertificateRequest{
							{
								CommonName:   "vvw-example.test",
								OU:           "VVW",
								Organisation: "VVW",
								Locality:     "London",
								Province:     "London",
								Country:      "GB",
								TTL:          "1h",
							},
						},
					},
				},
			},
		},
	},
}

const validPKIMonitorConfig = `
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

plugin "venafi-pki-monitor" "venafi-pki" {
  version = "v0.9.0"

  role "web_server" {
    secret "cloud" {
      venafi_cloud {
        apikey = "apikey"
      }
    }

	enforcement_policy {
	  zone = "zone"
	}

    import_policy {
      zone = "zone2"
    }

	intermediate_certificate {
      zone = "zone3"
	  common_name = "Vault SubCA"
	  ou = "VVW"
	  organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "1h"
	}

    allow_any_name = true
    ttl = "1h"
    max_ttl = "2h"

	test_certificate {
	  common_name = "vvw-example.test"
	  ou = "VVW"
	  organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "1h"
	}
  }
}`

var validPKIMonitorConfigResult = &Config{
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
			Type:      "venafi-pki-monitor",
			MountPath: "venafi-pki",
			Version:   "v0.9.0",
			Config:    nil,
			Impl: &pki_monitor.VenafiPKIMonitorConfig{
				MountPath: "venafi-pki",
				Version:   "v0.9.0",
				Role: pki_monitor.Role{
					Name:         "web_server",
					AllowAnyName: true,
					TTL:          "1h",
					MaxTTL:       "2h",
					EnforcementPolicy: &pki_monitor.Policy{
						Zone: "zone",
					},
					ImportPolicy: &pki_monitor.Policy{
						Zone: "zone2",
					},
					IntermediateCert: &pki_monitor.IntermediateCertRequest{
						Zone: "zone3",
						CertificateRequest: venafi.CertificateRequest{
							CommonName:   "Vault SubCA",
							OU:           "VVW",
							Organisation: "VVW",
							Locality:     "London",
							Province:     "London",
							Country:      "GB",
							TTL:          "1h",
						},
					},
					Secret: venafi.VenafiSecret{
						Name: "cloud",
						Cloud: &venafi.VenafiCloudConnection{
							APIKey: "apikey",
						},
					},
					TestCerts: []venafi.CertificateRequest{
						{
							CommonName:   "vvw-example.test",
							OU:           "VVW",
							Organisation: "VVW",
							Locality:     "London",
							Province:     "London",
							Country:      "GB",
							TTL:          "1h",
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
