---
layout: "venafi-vault-wizard"
page_title: "VVW: venafi-pki-monitor"
description: |-
Venafi Vault Wizard plugin that installs the venafi-pki-monitor Vault plugin to a Vault cluster.
---

# Plugin: venfai-pki-monitor

Configures the venafi-pki-monitor plugin to be installed into a Vault cluster through the Venafi Vault Wizard.


## Example Usage

The following example demonstrates the use of the venafi-pki-monitor plugin configuration.

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

$ vault write pki-backend/issue/tpp-backend common_name=test.test.test
```

* `secret` - (Required)
* `enforcement_policy` - (Optional)  
* `import_policy` - (Optional)
* `intermediate_certificate` - (Optional)
* `optional_config` - (Optional)  
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

#### enforcement_policy

* `zone` - (Required)

#### import_policy

* `zone` - (Required)

#### intermediate_certificate

* `common_name` - (Required)
* `ou` - (Required)
* `organisation` - (Required)
* `locality` - (Required)
* `province` - (Required)
* `country` - (Required)
* `ttl` - (Required)

#### optional_config
      
* `allow_any_name` - (Required)
* `ttl` - (Required)
* `max_ttl` - (Required)

### test_certificate

An optional test certificate to request, in order to verify everything is configured correctly.

* `common_name` - (Required)
* `ou` - (Required)
* `organisation` - (Required)
* `locality` - (Required)
* `province` - (Required)
* `country` - (Required)
* `ttl` - (Required)