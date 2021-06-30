# `venafi-vault-wizard`, (`vvw`)

-	Website: TBC
-	Source: https://github.com/opencredo/venafi-vault-wizard
-	Venafi Community Support: https://support.venafi.com/hc/en-us/community/topics
-	Venafi Cloud: https://www.venafi.com/venaficloud
-	HashiCorp Discuss: https://discuss.hashicorp.com/c/vault/30

This repository is home to the `venafi-vault-wizard` which can be used to verify the setup of HashiCorp Vault with Venafi Cloud and TPP.

## Requirements

-	[Go](https://golang.org/doc/install) 1.16
-	[Vagrant](https://www.vagrantup.com/downloads)
-	[Venafi Cloud Account & Zone](https://ui.venafi.cloud/login) or TPP instance

## Installation

While this tool is in development, you are required to build it yourself in order to use it.
To do this, simply run:

```shell
$ make build
```

Once this runs successfully, you can test it as follows:

```
$ ./bin/vvw 
VVW is a wizard to automate the installation and verification of Venafi PKI plugins for HashiCorp Vault.

Usage:
  vvw [command]

Available Commands:
  apply           Applies desired state as specified in config file
  generate-config Generates config file based on asking questions
  help            Help about any command

Flags:
  -f, --configFile string   Path to config file to use to configure Venafi Vault plugin (default "vvw_config.hcl")
  -h, --help                help for vvw

Use "vvw [command] --help" for more information about a command.
```

## Quick Start

### Single node Vault server

To quickly start exploring the use of the Venafi Vault Wizard, (VVW) a test environment with a VM running Vault can be easily set up using Vagrant.
This will provide the VVW tool with a Vault server to install the [vault-pki-backend-venafi](https://github.com/Venafi/vault-pki-backend-venafi) and [vault-pki-monitor-venafi](https://github.com/Venafi/vault-pki-monitor-venafi) plugins.
After they have been installed, certificates can be requested directly from the Vault instance.

First build the Venafi Vault Wizard, (VVW) tool. The binary will be placed in `./bin` at the root of the project.

```shell
$ make build
```

Navigate to the single-node test environment directory.
This directory contains a `Vagrantfile` and required scripts, as well as a sample `vvw.hcl` file to configure VVW appropriately. 
There is a `README.md` there which explains the setup in more detail.

```shell
$ cd examples/single_node_cluster_vagrant
```

### Multi-node Vault HA Cluster

Similarly, there is another test environment which does the same thing as the single node one above, with the difference that it provisions Vault in High Availability mode.
This starts three separate VMs which interact to form an HA cluster.
Again, there is a more detailed `README.md` in that directory.

```shell
$ cd examples/integrated_storage_ha_cluster_vagrant
```

Once the VVW tool has successfully completed the installation, a certificate can be requested from either plugin through Vault.
Replace `venafi-pki/issue/tls` with whatever mount path and role name was configured in the `vvw.hcl` file used.

```shell
$ vault write venafi-pki/issue/tls common_name="test.example.com"
```

## Configuration

The tool is configured using a configuration file written in HCL.
This supports two top level blocks, `vault {}` and `plugin "type" "mount path" {}`.

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

    allow_any_name = true
    ttl = "1h"
    max_ttl = "2h"

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

## Testing

The unit tests can be run with:

```shell
$ make test
```

### Generating mocks

The VVW tests use a number of pre-generated mocks that can be found under the `<repo root>/mocks` directory.
These replace the implementation of interfaces used throughout the code, to allow the tests to focus on testing specific areas.
They also provide the advantage that most unit tests run without touching real resources so are much faster and don't cause unwanted side effects.
If any of the interfaces have changed, or new ones added, then the mocks can be regenerated with the following command:

```shell
$ make generate-mocks
```

The command will download the [Mockery](http://github.com/vektra/mockery/v2@v2.6.0) binary to the `<repo root>/bin` directory and 
then proceed to generate mock implementations of interfaces found within the project.
See the [testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock) package for more details on how to use the mocking framework.
