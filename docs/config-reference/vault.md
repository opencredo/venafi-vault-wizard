# Configuration: vault

Venafi Vault Wizard configuration that represents a Vault cluster where Vault plugins will be installed and configured.

This must contain an `api_address` attribute pointing to the Vault server (this must be the leader node if running in HA mode), 
and a `token` attribute referencing a token with suitable permissions for configuring the plugins.

If installing the plugin on Vault servers with SSH access, and you would like the VVW tool to install the plugin binaries on the servers, the `vault`
block must then also contain an `ssh` block for each node in the cluster (only one if not running in HA mode).
If there are other nodes than those specified by the `ssh` blocks then the plugin won't be installed on all of them and strange behaviour may occur.
If Vault is running in containers, and the container image include the plugin binary already, then the `ssh` blocks can be omitted.

## Example Usage

The following example demonstrates the use of the `vault` configuration block to identify an existing HashiCorp Vault cluster.
There is an `env()` function available for pulling any of the values from environment variables.

```hcl
vault {
  api_address = "http://192.168.33.10:8200"
  token = env("VAULT_TOKEN")

  ssh {
    hostname = "192.168.33.10"
    port = 22
    username = "vagrant"
    password = "vagrant"
  }

  ssh {
    hostname = "192.168.33.11"
    port = 22
    username = "vagrant"
    password = "vagrant"
  }

  ssh {
    hostname = "192.168.33.12"
    port = 22
    username = "vagrant"
    password = "vagrant"
  }
}
```

## Argument Reference

The following arguments are supported:

* `api_address` - (Required) A string representing the `<protocol>://<host>:<port>` of the Vault cluster leader.
  This is the same value that you would set the `VAULT_ADDR` environment variable to for use with the `vault` CLI tool.
* `token` - (Required) A string representing a Vault token with enough privileges to install and configure Vault plugins.
* `ssh` - (Optional) A block representing location and credentials to used when access a node in the Vault cluster.

### SSH

The `ssh` block allows you to specify the location of a single node in a HashiCorp Vault Cluster.
The block can be specified multiple times depending on the size of the cluster.  

~> **Note:** The `ssh` blocks can be omitted if the plugin binaries are already installed to every node by external means.
For example if the Vault servers are running in containers, and the container images already have the plugins in them.

* `hostname` - (Required) A string representing the name or IP Address of a node in the Vault cluster.
* `port`     - (Required) A string representing the port that exposes SSH.
* `username` - (Required) A string representing the username of the SSH account on the node.
* `password` - (Required) A string representing the password of the SSH account on the node.
