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

plugin "venafi-pki-backend" "pki-backend" {
  role "cloud" {
    secret "cloud" {
      venafi_cloud {
        apikey = env("VENAFI_API_KEY")
        zone = "6225eee0-8101-11eb-7822-0b1983e1b167"
      }
    }

    test_certificate {
      common_name = "vvw-example.test"
    }
  }
}

