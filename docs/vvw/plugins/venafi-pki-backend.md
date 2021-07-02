---
layout: "venafi-vault-wizard"
page_title: "VVW: venafi-pki-backend"
description: |-
Venafi Vault Wizard plugin that installs the venafi-pki-backend Vault plugin to a Vault cluster.
---

# Plugin: venafi-pki-backend

Configures the `venafi-pki-backend` plugin to be installed into a Vault cluster through the Venafi Vault Wizard. 


## Example Usage

The following example demonstrates the use of the `venafi-pki-backend` plugin configuration.
The `version` argument is common to all plugins and is described in [the plugin block's documentation](../plugin.md).

```hcl
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

## Argument Reference

The following arguments are supported:

* `role` - (Required) A block corresponding to a role within the plugin, from which certificates can be requested.

### role

A `role` block is given a label that specifies the role name.
This is what will be used in the path, along with the mount path, to request certificates.

```hcl
role "tpp-backend" {
  ...
}

$ vault write pki-backend/issue/tpp-backend common_name=test.test.test
```

The `role` block supports the following blocks:

* `secret` - (Required) The Venafi connection details
* `test_certificate` - (Optional) Details to use to request a certificate from the role, in order to verify that it is configured correctly.

#### secret

The `secret` block must contain exactly one of the following blocks:

* `venafi_tpp` - Required when using Venafi's Trust Protection Platform
* `venafi_cloud` - Required when using Venafi as a Service

##### venafi_tpp

* `url` - (Required)  A String specifying the URL of the API endpoint for the Venafi Trust Protection Platform, (TPP).
* `username` - (Required) A string representing a TPP account username
* `password` - (Required) A string representing a TPP account password
  Avoid hardcoding this in the configuration file in case it gets leaked.
  It is recommended to use `env("TPP_PASSWORD")` to retrieve this from an environment variable instead.
* `zone` - (Required) The path of the policy within TPP, from which certificates will be requested.

##### venafi_cloud

* `apikey` - (Required) A String with an API Key with access to Venafi as a Service.
  Avoid hardcoding this in the configuration file in case it gets leaked.
  It is recommended to use `env("VENAFI_API_KEY")` to retrieve this from an environment variable instead.
* `zone` - (Required) The project zone within Venafi as a Service from which certificates will be requested.
  The format is `Application Name\Issuing Template API Alias`.

### test_certificate

An optional test certificate to request, in order to verify everything is configured correctly.
The arguments correspond to the usual parameters found in a certificate signing request (CSR).

* `common_name` - (Required)
* `ou` - (Required)
* `organisation` - (Required)
* `locality` - (Required)
* `province` - (Required)
* `country` - (Required)
* `ttl` - (Required)