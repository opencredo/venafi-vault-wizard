---
layout: "venafi-vault-wizard"
page_title: "VVW: venafi-pki-backend"
description: |-
Venafi Vault Wizard plugin that installs the venafi-pki-backend Vault plugin to a Vault cluster.
---

# Plugin: venfai-pki-backend

Configures the venafi-pki-backend plugin to be installed into a Vault cluster through the Venafi Vault Wizard. 


## Example Usage

The following example demonstrates the use of the venafi-pki-backend plugin configuration.

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

* `version` - (Required)

* `role` - (Required)

* `test_certificate` - (Optional) A test certificate request used to verify the plugin is configured correctly.

### role

A role block is given a label that is then used to configure the plugins path

```hcl
role "tpp-backend" {
  ...
}

vault write pki-backend/issue/tpp-backend common_name=test.test.test
```

* `secret` - (Required)
* `test_certificate` - (Optional)

#### secret

* `venafi_tpp` - (Optional/Required for TPP backends)
* `venafi_cloud` - (Optional/Required for Venafi Cloud)

##### venafi_tpp

* `url` - (Required)  A String representing the URL endpoint for the Venafi Trust Protection Platform, (TPP).
* `username` - (Required) A string representing a TPP account username
* `password` - (Required) A string representing a TPP account password
* `zone` - (Required)

##### venafi_cloud

* `apikey` - (Required) A string repsenting a Venafi Cloud generated APK Key.
* `zone` - (Required)

### test_certificate