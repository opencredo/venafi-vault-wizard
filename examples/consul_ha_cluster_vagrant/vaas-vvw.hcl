vault {
  api_address = "http://192.168.56.10:8200"
  token = env("VAULT_TOKEN")

  ssh {
    hostname = "192.168.56.10"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }

  ssh {
    hostname = "192.168.56.11"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }

  ssh {
    hostname = "192.168.56.12"
    username = "vagrant"
    password = "vagrant"
    port = 22
  }
}

plugin "venafi-pki-monitor" "pki-monitor" {
  version = "v0.9.0"

  # A role called "web_server" can be used with:
  # vault write pki-monitor/issue/web_server common_name=test.test.test
  role "web_server" {

    # Connection details for Venafi VaaS
    # If using Venafi TPP, replace the venafi_vaas block with a venafi_tpp one and specify the "url", "username" and "password" attributes instead
    secret "vaas" {
      venafi_vaas {
        apikey = env("VENAFI_API_KEY")
      }
    }

    # Policy to use to specify rules for issuing certificates
    enforcement_policy {
      zone = "VVW Test\\VVW SubCA"
    }

    # Details of the root certificate with which to issue certificates
    root_certificate {
      common_name = "Vault SubCA"
      ou = "OpenCredo"
      organisation = "VVW"
      locality = "London"
      province = "London"
      country = "GB"
      ttl = "3h"
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
  version = "v0.10.3"

  # A role called "web_server" can be used with:
  # vault write pki-backend/issue/web_server common_name=test.test.test
  role "web_server" {

    # Connection details for Venafi VaaS
    # If using Venafi TPP, replace the venafi_vaas block with a venafi_tpp one and specify the "url", "username" and "password" attributes instead
    secret "vaas" {
      zone = "VVW Test\\VVW SubCA"
      venafi_vaas {
        apikey = env("VENAFI_API_KEY")
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
