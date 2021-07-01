# Configuration File Syntax

The configuration file supports two top level blocks, `vault {}` and `plugin "type" "mount path" {}`.
The format of each block is described in the following sections.

## `vault` block

An example of the `vault` block (taken from the `test_envs/integrated_storage_ha_cluster_vagrant` test environment) is shown below:

```hcl
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
```

This must contain an `api_address` attribute pointing to the Vault server (this must be the leader node if running in HA mode), and a `token` attribute referencing a token with suitable permissions for configuring the plugins.
There is an `env()` function available for pulling any of the values from environment variables.

The `vault` block must then also contain an `ssh` block for each node in the cluster (only one if not running in HA mode).
If there are other nodes than those specified by the `ssh` blocks then the plugin won't be installed on all of them and strange behaviour may occur.
The `ssh` blocks can be omitted if the plugin binaries are already installed to every node by external means.
For example if the Vault servers are running in containers, and the container images already have the plugins in them.


## `plugin` block

The structure of the `plugin` block depends on the specific plugin being used.
The first block label should specify which plugin to install (currently only `venafi-pki-backend` and `venafi-pki-monitor` are supported).
The second block label specifies which path the plugin will be mounted at in Vault.
Every plugin must also specify a `version` attribute too.

For the two currently supported plugins, use the following examples as templates to modify for your own needs:

```hcl
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

    # Details of the intermediate certificate with which to issue certificates
    intermediate_certificate {
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
      ttl = "1h"
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
  }
}
```