vault {
  # No ssh blocks needed as the plugin binaries are already installed to the containers
  api_address = "http://localhost:8200"
  token = env("VAULT_TOKEN")
}

plugin "venafi-pki-monitor" "pki-monitor" {
  version = "v0.9.0"

  # A role called "web_server" can be used with:
  # vault write pki-monitor/issue/web_server common_name=test.test.test
  role "web_server" {
    # Connection details for Venafi TPP
    # If using Venafi Cloud, replace the venafi_tpp block with a venafi_cloud one and specify the "apikey" attribute instead
    secret "tpp" {
      venafi_tpp {
        url = env("TPP_URL")
        username = env("TPP_USERNAME")
        password = env("TPP_PASSWORD")
      }
    }

    # Policy to use to specify rules for issuing certificates
    enforcement_policy {
      zone = "Partner Dev\\TLS\\Certificates\\HashiCorp Vault\\Vault SubCA"
    }

    # Policy to send details of issued certificates to
    import_policy {
      zone = "Partner Dev\\TLS\\Certificates\\HashiCorp Vault\\Vault Issued"
    }

    # Details of the root certificate with which to issue certificates
    root_certificate {
      common_name = "Vault SubCA"
      ou = "OpenCredo"
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

    # An optional test certificate to request, in order to verify everything works
    test_certificate {
      common_name = "test1.venafidemo.com"
      ou = "OpenCredo"
      organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "5m"
    }

    test_certificate {
      common_name = "test2.venafidemo.com"
      ou = "OpenCredo"
      organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "5m"
    }
  }
}

plugin "venafi-pki-backend" "pki-backend" {
  version = "v0.9.0"

  # A role called "tpp-backend" can be used with:
  # vault write pki-backend/issue/tpp-backend common_name=test.test.test
  role "tpp-backend" {
    # Connection details for Venafi TPP
    # If using Venafi Cloud, replace the venafi_tpp block with a venafi_cloud one and specify the "apikey" attribute instead
    secret "tpp" {
      venafi_tpp {
        url = env("TPP_URL")
        username = env("TPP_USERNAME")
        password = env("TPP_PASSWORD")
        zone = "Partner Dev\\TLS\\Certificates\\HashiCorp Vault\\Vault PKI Backend"
      }
    }

    # An optional test certificate to request, in order to verify everything works
    test_certificate {
      common_name = "test1.venafidemo.com"
      ou = "OpenCredo"
      organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "1h"
    }

    test_certificate {
      common_name = "test2.venafidemo.com"
      ou = "OpenCredo"
      organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "1h"
    }
  }
}
