# `venafi-vault-wizard`, (`vvw`)

-	Website: TBC
-	Source: https://github.com/opencredo/venafi-vault-wizard
-	Venafi Community Support: https://support.venafi.com/hc/en-us/community/topics
-	Venafi Cloud: https://www.venafi.com/venaficloud
-	HashiCorp Discuss: https://discuss.hashicorp.com/c/vault/30

This repository is home to the `venafi-vault-wizard` which can be used to verify the setup of HashiCorp Vault with Venafi-as-a-Service and TPP.

## Requirements

-	[Go](https://golang.org/doc/install) 1.16
-	[Vagrant](https://www.vagrantup.com/downloads)
-	[Venafi Cloud Account & Zone](https://ui.venafi.cloud/login) or TPP instance

## Docs

- [Installation](docs/installation.md)
- [Config File Syntax](docs/config-file-format.md)
- [Generating Config Files with the Step-by-Step Wizard](docs/config-generation.md)
- [Example Environments](examples/README.md)
- [VVW Supported Plugins](docs/vvw/index.md)

## Introduction

The tool centres around the use of a configuration file to declare what plugins should be installed to HashiCorp Vault, and how they should be configured.
This file describes the desired state, and the tool uses it to make the required changes to achieve the desired state.
This workflow is similar to that of HashiCorp Terraform.

The tool has two main subcommands: `generate-config` and `apply`.
The `generate-config` command triggers a step-by-step wizard that asks a series of questions in order to generate the configuration file.
The `apply` command, similarly to Terraform, then performs the plugin installation and configuration, as required by the configuration file.
Both subcommands have a required `-f` or `--configFile` flag.
For `generate-config`, this specifies where the generated configuration will be written to.
For `apply`, it specifies the configuration to read from, and to apply to the Vault server.

## Quick Start

### Single node Vault server

To quickly start exploring the use of the Venafi Vault Wizard, (VVW) a test environment with a VM running Vault can be easily set up using Vagrant.
This will provide the VVW tool with a Vault server to install the [vault-pki-backend-venafi](https://github.com/Venafi/vault-pki-backend-venafi) and [vault-pki-monitor-venafi](https://github.com/Venafi/vault-pki-monitor-venafi) plugins.
After they have been installed, certificates can be requested directly from the Vault instance.
A Venafi TPP instance must be available, or alternatively Venafi-as-a-Service can be used with some minor modifications.

First, build the Venafi Vault Wizard, (VVW) tool. The binary will be placed in `./bin` at the root of the project.

```shell
$ make build
```

Navigate to the single-node test environment directory.
This directory contains a `Vagrantfile` and required scripts, as well as a sample `vvw.hcl` file to configure VVW appropriately. 
There is a `README.md` there which explains the setup in more detail.

```shell
$ cd examples/single_node_cluster_vagrant
```

Provision the test Vault server and set the required environment variables using the following commands, substituting the relevant information where appropriate.
If using Venafi-as-a-Service, some tweaks will need to be made to the configuration file to remove the references to TPP.
See the [configuration file documentation](docs/config-file-format.md) for more information on this.

```shell
$ vagrant up
$ export VAULT_TOKEN="TOKEN PRINTED FROM VAGRANT HERE"
$ export VAULT_ADDR="http://192.168.33.20:8200"
$ export TPP_URL="YOUR TPP INSTANCE URL HERE"
$ export TPP_USERNAME="YOUR TPP USERNAME HERE"
$ export TPP_PASSWORD="YOUR TPP PASSWORD HERE"
$ export VENAFI_API_KEY="YOUR VaaS API KEY"
```

When that has finished, run the VVW tool with the provided `vvw.hcl` configuration file:
There are two vvw HCL configuration files.  `tpp-vvw.hcl` for Trust Protection Platform and `vaas-vvw.hcl` for Venafi as a Service.

```shell
$ ../../bin/vvw apply -f tpp-vvw.hcl
```
or for the Vass configuration
```shell
$ ../../bin/vvw apply -f vaas-vvw.hcl
```

Once the VVW tool has successfully completed the installation, a certificate can be requested from either plugin through Vault.

```shell
$ vault write pki-monitor/issue/web_server common_name="test.example.com"
$ vault write pki-backend/issue/tpp-backend common_name="test.example.com"
```
or for the VaaS configuration
```shell
$ vault write pki-monitor/issue/web_server common_name="test.example.com"
$ vault write pki-backend/issue/vaas-backend common_name="test.example.com"
```

## Development

### Testing

The unit tests can be run with:

```shell
$ make test
```

#### Generating mocks

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
