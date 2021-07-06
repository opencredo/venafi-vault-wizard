---
layout: "venafi-vault-wizard"
page_title: "VVW: plugin"
description: |-
Venafi Vault Wizard configuration that represents a plugin to be installed and configured.
---

# Configuration: plugin

Venafi Vault Wizard configuration that represents a plugin to be installed and configured.

There are some options that are relevant to every plugin, including the mount path, and the plugin version.
Each plugin then provides its own specific configuration options that are discussed in their own plugin-specific documentation page.

There can be multiple `plugin` blocks in one configuration file, and these each correspond to one plugin mount in Vault.
Multiple instances of the same plugin can be mounted multiple times, as long as they have unique mount paths.

The `plugin` block is always followed by two labels.
The first is the plugin type, i.e. which plugin should be installed.
The second label is the mount path, or the Vault path under which all the configuration can be accessed.

## Example Usage

The following example demonstrates the use of the `plugin` configuration block to install and configure version 0.9.0 of the `venafi-pki-backend` plugin.
This example specifies that it should be mounted at `pki-backend/` in Vault, which will mean that the CLI commands to access it will look like `vault write pki-backend/...`.
The `role` block in this example is specific to that plugin, so the documentation for that can be found in the `venafi-pki-backend` page.

```hcl
plugin "venafi-pki-backend" "pki-backend" {
  version = "v0.9.0"

  # A role called "web_server" can be used with:
  # vault write pki-backend/issue/web_server common_name=test.test.test
  role "web_server" {
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
  }
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Required) The version of the plugin that should be installed.
  The format will vary depending on which plugin is being installed, but will usually be the Git tag/release name of the plugin.
* `filename` - (Optional) The filename of the plugin binary on the Vault server filesystem.
  This will default to `pluginType_version` if not specified, and is only recommended for use if the plugin binaries are installed by external means, and the filename can't be changed.
* `build_arch` - (Optional) The OS and CPU architecture of the Vault server.
  Defaults to `linux`.
  Options are: `linux`, `linux86`, `darwin`, `windows`, `windows86`.
