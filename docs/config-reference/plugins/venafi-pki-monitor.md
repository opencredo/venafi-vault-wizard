# Plugin: venfai-pki-monitor

Configures the `venafi-pki-monitor` plugin to be installed into a Vault cluster through the Venafi Vault Wizard.


## Example Usage

The following example demonstrates the use of the `venafi-pki-monitor` plugin configuration.
The `version` argument is common to all plugins and is described in [the plugin block's documentation](../plugin.md).

```hcl
plugin "venafi-pki-monitor" "pki-monitor" {
  version = "v0.9.0"

  # A role called "web_server" can be used with:
  # vault write pki-monitor/issue/web_server common_name=test.test.test
  role "web_server" {
    # Connection details for Venafi TPP
    # If using Venafi as a Service, replace the venafi_tpp block with a venafi_vaas one and specify the "apikey" attribute instead
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

* `role` - (Required) A block corresponding to a role within the plugin, from which certificates can be requested.

### role

A `role` block is given a label that specifies the role name.
This is what will be used in the path, along with the mount path, to request certificates.

```hcl
role "web_server" {
  ...
}

$ vault write pki-backend/issue/web_server common_name=test.test.test
```

The `role` block supports the following blocks:

* `secret` - (Required)
* `enforcement_policy` - (Optional)  
* `import_policy` - (Optional)
* `intermediate_certificate` - (Optional)
* `optional_config` - (Optional)  
* `test_certificate` - (Optional)

#### secret

The `secret` block must contain exactly one of the following blocks:

* `venafi_tpp` - Required when using Venafi's Trust Protection Platform
* `venafi_vaas` - Required when using Venafi as a Service

##### venafi_tpp

* `url` - (Required)  A String representing the URL endpoint for the Venafi Trust Protection Platform, (TPP).
* `username` - (Required) A string representing a TPP account username
* `password` - (Required) A string representing a TPP account password

~> **Warning:** Avoid hardcoding this in the configuration file in case it gets leaked.
It is recommended to use `env("TPP_PASSWORD")` to retrieve this from an environment variable instead.

##### venafi_vaas

* `apikey` - (Required) A String with an API Key with access to Venafi as a Service.

~> **Warning:** Avoid hardcoding this in the configuration file in case it gets leaked.
It is recommended to use `env("VENAFI_API_KEY")` to retrieve this from an environment variable instead.

#### enforcement_policy

* `zone` - (Required)

#### import_policy

* `zone` - (Required)

#### intermediate_certificate

* `common_name` - (Required) The fully qualified domain name (FQDN) of your server.  For example `www.example.com`
* `ou` - (Required) The legal name of your organization.
* `organisation` - (Required) The division of your organization handling the certificate.
* `locality` - (Required) The city where your organization is located.
* `province` - (Required) The state/region where your organization is located
* `country` - (Required) The two-letter code for the country where your organization is located.
* `ttl` - (Required) The Time To Live for your certificate

#### optional_config
      
* `allow_any_name` - (Required)
* `ttl` - (Required)
* `max_ttl` - (Required)

### test_certificate

An optional test certificate to request, in order to verify everything is configured correctly.
The arguments correspond to the usual parameters found in a certificate signing request (CSR).

* `common_name` - (Required) The fully qualified domain name (FQDN) of your server.  For example `www.example.com`
* `ou` - (Required) The legal name of your organization.
* `organisation` - (Required) The division of your organization handling the certificate.
* `locality` - (Required) The city where your organization is located.
* `province` - (Required) The state/region where your organization is located
* `country` - (Required) The two-letter code for the country where your organization is located.
* `ttl` - (Required) The Time To Live for your certificate