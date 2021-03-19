# venafi-vault-wizard, (vvw)

-	Website: TBC
-	Source: https://github.com/opencredo/venafi-vault-wizard
-	Venafi Community Support: https://support.venafi.com/hc/en-us/community/topics
-	Venafi Cloud: https://www.venafi.com/venaficloud
-	HashiCorp Discuss: https://discuss.hashicorp.com/c/vault/30

This repository is home to the `venafi-vault-wizard` which can be used to verify the setup of HashiCorp Vault and Venafi Cloud.

## Requirements

-	[Go](https://golang.org/doc/install) 1.16
-	[Vagrant](https://www.vagrantup.com/downloads)
-	[Venafi Cloud Account & Zone](https://ui.venafi.cloud/login)

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
  help        Help about any command
  install     Installs a Venafi plugin to Vault

Flags:
  -h, --help                    help for vvw
      --sshPassword string      Password for SSH user to log into Vault server with (default "password")
      --sshPort uint            Port on which SSH is running on the Vault server (default 22)
      --sshUser string          Username with which to log into Vault server over SSH (must have sudo privileges) (default "username")
      --vaultAddress string     Vault HTTP API endpoint (default "https://127.0.0.1:8200")
      --vaultMountPath string   Vault path at which to mount the Venafi plugin (default "venafi-pki")
      --vaultToken string       Token used to authenticate with Vault (default "root")
      --venafiAPIKey string     API Key used to access Venafi Cloud
      --venafiZone string       Venafi Cloud Project Zone in which to create certificates

Use "vvw [command] --help" for more information about a command.
```

## Quick Start

To quickly start exploring the use of the Venafi Vault Wizard, (VVW) a demo environment 
with a VM running Vault can be easily set up using Vagrant.  This will provide the VVW tool with a 
Vault server to install the [vault-pki-backend-venafi](https://github.com/Venafi/vault-pki-backend-venafi) 
and [vault-pki-monitor-venafi](https://github.com/Venafi/vault-pki-monitor-venafi) plugins.  After the VVW has installed the 
Venafi Vault plugins certificates can be requested.

First build the Venafi Vault Wizard, (VVW) tool. The binary will be placed in `./bin` at the root of the project.

```shell
$ make build
```

Navigate to the demo environment directory. This directory contains Vagrantfile and supportive scripts. 
Start the demo environment through the `vagrant up` command.

```shell
$ cd test_envs/vagrant
$ vagrant up
```

Set your environment variables for Vault and Venafi.  Vagrant will return the Vault token as the last line of its output.

The [Venafi Cloud UI](https://ui.venafi.cloud/login) locate your API Key and Zone ID.
Once setup a `vault status` can be used to inspect the deployed Vault server.

By default, a host-only network is created, with Vault at 192.168.33.10.
This can be changed in the `Vagrantfile` if needed.

```shell
$ export VAULT_ADDR=http://192.168.33.10:8200
$ export VAULT_TOKEN=<VAULT TOKEN HERE>
$ export VENAFI_CLOUD_API_KEY=<VENAFI CLOUD API KEY HERE>
$ export VENAFI_CLOUD_ZONE_ID=<VENAFI CLOUD ZONE ID HERE>
$ vault status
```

The following command will execute the VVW tool against the Vault server and install the install `venafi-pki-backend` Vault Plugin.
The VVW tool will provide a progress report as the installation progresses.

```shell
$ ../../bin/vvw install venafi-pki-backend \
  --vaultAddress=$VAULT_ADDR \
  --vaultToken=$VAULT_TOKEN \
  --venafiAPIKey=$VENAFI_CLOUD_API_KEY \
  --venafiZone=$VENAFI_CLOUD_ZONE_ID \
  --sshUser=vagrant \
  --sshPassword=vagrant \
  --sshPort=22
```

Once the VVW tool has successfully completed a certificate can be requested through Vault.

```shell
$ vault write venafi-pki/issue/cloud common_name="test.example.com"
```

## Generating Test Mocks

The VVW tests use a number of pre-generated mocks that can be found under the `<repo root>/mocks` directory and allow the 
tests to be executed upon checkout.  To generate new mocks the following command can be used.

```shell
$ make generate-mocks
```

The command will download the [Mockery](http://github.com/vektra/mockery/v2@v2.6.0) binary to the `<repo root>/bin` directory and 
then proceed to generate mock implementations of interfaces found within the project.
