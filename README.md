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
$ cd test_envs/vagrant
$ vagrant up
$ export VAULT_ADDR=http://192.168.33.10:8200
$ vault status
$ ./vvw install venafi-pki-backend \
  --vaultAddress=http://192.168.33.10:8200 \
  --sshAddress=192.168.33.10:22 \
  --vaultToken="<DISPLAYED IN VAGRANT UP OUTPUT>" \
  --venafiAPIKey="<VENAFI CLOUD API KEY HERE>" \
  --venafiZone="<VENAFI CLOUD ZONE ID HERE>"
$ vault write venafi-pki/issue/cloud common_name="test.example.com"
```
