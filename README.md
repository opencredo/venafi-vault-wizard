# `venafi-vault-wizard`
This repository is home to the `venafi-vault-wizard` which can be used to verify the setup of HashiCorp Vault and Venafi Cloud.

## Requirements

-	[Go](https://golang.org/doc/install) 1.16

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

## Testing with Vagrant

A demo environment with a VM running Vault can be easily set up using Vagrant.
This is to allow testing of the VVW tool.
By default, a host-only network is created, with Vault at 192.168.33.10.
This can be changed in the `Vagrantfile` if needed.

```shell
$ cd test_envs/vagrant
$ vagrant up
$ export VAULT_ADDR=http://192.168.33.10:8200
$ export VAULT_TOKEN="<WHATEVER WAS SHOWN IN VAGRANT LOGS>"
$ vault status
$ ./bin/vvw install venafi-pki-backend \
  --venafiAPIKey="<VENAFI CLOUD API KEY HERE>" \
  --venafiZone="<VENAFI CLOUD ZONE ID HERE>"
  --sshUser=vagrant \
  --sshPassword=vagrant

$ vault write venafi-pki/issue/cloud common_name="test.example.com"
```
