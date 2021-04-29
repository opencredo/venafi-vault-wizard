vault {
  api_address = "http://192.168.33.10:8200"
  token = env("VAULT_TOKEN")

  ssh {
    hostname = "192.168.33.10"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }

  ssh {
    hostname = "192.168.33.11"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }

  ssh {
    hostname = "192.168.33.12"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-monitor" "pki-monitor" {
  version = "v0.9.0"

  role "web_server" {
    secret "tpp" {
      venafi_tpp {
        url = env("TPP_URL")
        username = env("TPP_USERNAME")
        password = env("TPP_PASSWORD")
      }
    }

    enforcement_policy {
      zone = "Partner Dev\\TLS\\Certificates\\HashiCorp Vault\\Vault SubCA"
    }

    import_policy {
      zone = "Partner Dev\\TLS\\Certificates\\HashiCorp Vault\\Vault Issued"
    }

    intermediate_certificate {
      common_name = "Vault SubCA"
      ou = "OpenCredo"
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

plugin "venafi-pki-backend" "pki-backend" {
  version = "v0.9.0"

  role "tpp-backend" {
    secret "tpp" {
      venafi_tpp {
        url = env("TPP_URL")
        username = env("TPP_USERNAME")
        password = env("TPP_PASSWORD")
        zone = "Partner Dev\\TLS\\Certificates\\HashiCorp Vault\\Vault PKI Backend"
      }
    }

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
