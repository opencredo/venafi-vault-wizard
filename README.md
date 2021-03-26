# `venafi-vault-wizard`
This repository is home to the `venafi-vault-wizard` which can be used to verify the setup of HashiCorp Vault and Venafi Cloud.

## Requirements

-	[Go](https://golang.org/doc/install) 1.16

## Testing with Vagrant

A demo environment with a VM running Vault can be easily set up using Vagrant.
This is to allow testing of the VVW tool.
By default, a host-only network is created, with Vault at 192.168.33.10.
This can be changed in the `Vagrantfile` if needed.

```shell
c cd test_envs/vagrant
$ vagrant up
$ export VAULT_ADDR=http://192.168.33.10:8200
$ vault status
$ ./vvw install venafi-pki-backend \
  --vaultAddress=http://192.168.33.10:8200 \
  --vaultToken=$VAULT_TOKEN \
  --venafiAPIKey="<VENAFI CLOUD API KEY HERE>" \
  --venafiZone="<VENAFI CLOUD ZONE ID HERE>"
  --sshUser=vagrant \
  --sshPassword=vagrant \
  --sshPort=22

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


