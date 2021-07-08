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
		"valid venafi-pki-backend with VaaS": {
			config:  validPKIBackendVaaSConfig,
			want:    validPKIBackendVaaSConfigResult,
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
		"invalid venafi-pki-backend with both VaaS and TPP": {
			config:  invalidPKIBackendConfigBothConnectionTypes,
			wantErr: true,
		},
		"valid venafi-pki-monitor with VaaS": {
			config:  validPKIMonitorConfig,
			want:    validPKIMonitorConfigResult,
			wantErr: false,
		},
		"invalid venafi-pki-monitor with blank intermediate certificate zone": {
			config:  invalidPKIMonitorConfigBlankIntermediateCertificateZone,
			wantErr: true,
		},
		"invalid venafi-pki-monitor with blank enforcement policy zone": {
			config:  invalidPKIMonitorConfigBlankEnforcementPolicyZone,
			wantErr: true,
		},
		"invalid venafi-pki-monitor with blank import policy zone": {
			config:  invalidPKIMonitorConfigBlankImportPolicyZone,
			wantErr: true,
		},
		"valid venafi-pki-backend with defined build arch": {
			config:  validPKIBackendBuildArchConfig,
			want:    validPKIBackendBuildArchConfigResult,
			wantErr: false,
		},
		"valid venafi-pki-monitor with defined build arch": {
			config:  validPKIMonitorBuildArchConfig,
			want:    validPKIMonitorBuildArchConfigResult,
			wantErr: false,
		},
		"invalid venafi-pki-backend with incorrect build arch": {
			config:  invalidPKIBackendBuildArchConfig,
			wantErr: true,
		},
		"invalid venafi-pki-monitor with incorrect build arch": {
			config:  invalidPKIMonitorBuildArchConfig,
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

const validPKIBackendVaaSConfig = `
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
  role "vaas" {
    secret "vaas" {
      zone = "zone1"
      venafi_vaas {
        apikey = "apikey"
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

var validPKIBackendVaaSConfigResult = &Config{
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
	Plugins: []plugins.PluginConfig{
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
						Name: "vaas",
						Secret: pki_backend.ZonedSecret{
							Name: "vaas",
							Zone:  "zone1",
							VenafiSecret: venafi.VenafiSecret{
								VaaS: &venafi.VenafiVaaSConnection{
									APIKey: "apikey",
								},
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
    secret "tpptest" {
      zone = "zone1"
      venafi_tpp {
        url = "tpp.venafitest.com"
        username = "admin"
        password = "pword234"
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
	Plugins: []plugins.PluginConfig{
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
						Secret: pki_backend.ZonedSecret{
							Name: "tpptest",
							Zone:     "zone1",
							VenafiSecret: venafi.VenafiSecret{
								TPP: &venafi.VenafiTPPConnection{
									URL:      "tpp.venafitest.com",
									Username: "admin",
									Password: "pword234",
								},
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
    secret "vaas" {
      venafi_vaas {
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

	optional_config {
      allow_any_name = true
      ttl = "1h"
      max_ttl = "2h"
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
	Plugins: []plugins.PluginConfig{
		{
			Type:      "venafi-pki-monitor",
			MountPath: "venafi-pki",
			Version:   "v0.9.0",
			Config:    nil,
			Impl: &pki_monitor.VenafiPKIMonitorConfig{
				MountPath: "venafi-pki",
				Version:   "v0.9.0",
				Role: pki_monitor.Role{
					Name: "web_server",
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
					Secret: pki_monitor.UnZonedSecret{
						Name: "vaas",
						VenafiSecret: venafi.VenafiSecret{
							VaaS: &venafi.VenafiVaaSConnection{
								APIKey: "apikey",
							},
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
					OptionalConfig: &venafi.OptionalConfig{
						AllowAnyName: true,
						TTL:          "1h",
						MaxTTL:       "2h",
					},
				},
			},
		},
	},
}

const validPKIBackendBuildArchConfig = `
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
  build_arch = "linux86"
  role "vaas" {
    secret "vaas" {
      zone = "zone1"
      venafi_vaas {
        apikey = "apikey"
      }
    }
  }
}`

var validPKIBackendBuildArchConfigResult = &Config{
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
	Plugins: []plugins.PluginConfig{
		{
			Type:      "venafi-pki-backend",
			MountPath: "venafi-pki",
			Version:   "v0.9.0",
			BuildArch: "linux86",
			Config:    nil,
			Impl: &pki_backend.VenafiPKIBackendConfig{
				MountPath: "venafi-pki",
				Version:   "v0.9.0",
				BuildArch: "linux86",
				Roles: []pki_backend.Role{
					{
						Name: "vaas",
						Secret: pki_backend.ZonedSecret{
							Name: "vaas",
							Zone: "zone1",
							VenafiSecret: venafi.VenafiSecret{
								VaaS: &venafi.VenafiVaaSConnection{
									APIKey: "apikey",
								},
							},
						},
					},
				},
			},
		},
	},
}

const validPKIMonitorBuildArchConfig = `
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
  build_arch = "linux86"

  role "web_server" {
    secret "vaas" {
      venafi_vaas {
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
  }
}`

var validPKIMonitorBuildArchConfigResult = &Config{
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
	Plugins: []plugins.PluginConfig{
		{
			Type:      "venafi-pki-monitor",
			MountPath: "venafi-pki",
			Version:   "v0.9.0",
			BuildArch: "linux86",
			Config:    nil,
			Impl: &pki_monitor.VenafiPKIMonitorConfig{
				MountPath: "venafi-pki",
				Version:   "v0.9.0",
				BuildArch: "linux86",
				Role: pki_monitor.Role{
					Name: "web_server",
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
					Secret: pki_monitor.UnZonedSecret{
						Name: "vaas",
						VenafiSecret: venafi.VenafiSecret{
							VaaS: &venafi.VenafiVaaSConnection{
								APIKey: "apikey",
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
  role "vaas" {
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

  role "vaas" {
    secret "vaas" {
      zone = "zone"
      venafi_vaas {
        apikey = "apikey"
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

const invalidPKIMonitorConfigBlankIntermediateCertificateZone = `
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
    secret "vaas" {
      venafi_vaas {
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
      zone = ""
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
  }
}`

const invalidPKIMonitorConfigBlankEnforcementPolicyZone = `
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
    secret "vaas" {
      venafi_vaas {
        apikey = "apikey"
      }
    }

	enforcement_policy {
	  zone = ""
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
  }
}`

const invalidPKIMonitorConfigBlankImportPolicyZone = `
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
    secret "vaas" {
      venafi_vaas {
        apikey = "apikey"
      }
    }

	enforcement_policy {
	  zone = "zone"
	}

    import_policy {
      zone = ""
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
  }
}`

const invalidPKIBackendBuildArchConfig = `
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
  build_arch = "linux386"
  role "vaas" {
    secret "vaas" {
      zone = "zone1"
      venafi_vaas {
        apikey = "apikey"
      }
    }
  }
}`

const invalidPKIMonitorBuildArchConfig = `
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
  build_arch = "linux386"

  role "web_server" {
    secret "vaas" {
      venafi_vaas {
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
  }
}`

func deletePluginsUncheckedFields(config *Config) {
	for i := 0; i < len(config.Plugins); i++ {
		config.Plugins[i].Config = nil
	}
}
