# Vault Test Environment with Single Node

This test environment spins up a VM running a single Vault node.
This node has the SSH credentials `vagrant` for both the username and password.
The node's IP addresses is `192.168.33.20`

## Setup

To use this `vvw.hcl` file as is, the following environment variables need to be set:

- `VAULT_TOKEN`: this will be printed to stdout inside the VM so should be reflected in the Vagrant output
- `TPP_URL`: this should be the API address ending in `/vedsdk` of the Venafi TPP instance
- `TPP_USERNAME`: this should be a TPP user suitable for use by the plugins
- `TPP_PASSWORD`: the corresponding password of the user above

Furthermore, the TPP environment needs to have a number of policies configured.
Specifically, the `venafi-pki-backend` plugin needs one policy, and the `venafi-pki-monitor` needs two.
Consult these blog posts for instructions on configuring TPP for use with the plugins:

- [`venafi-pki-monitor`](https://medium.com/hashicorp-engineering/vault-integration-patterns-with-venafi-21c3626cdcdb)
- [`venafi-pki-backend`](https://medium.com/hashicorp-engineering/vault-integration-patterns-with-venafi-part-2-ff6a5fcc3d3d)

Once these policies are set up, change the `zone` references in `vvw.hcl` to reference them.
Similarly, feel free to change the `intermediate_certificate` and `test_certificate` subject information to suit requirements, but ensure they comply with the policies configured in TPP.

If using Venafi Cloud, then adjust the `zone`s and `secret`s accordingly.

## Usage

To spin up the cluster, run:

```shell
$ vagrant up
```

This will take a little while as it will install Vault and configure things as required.
When this has completed, copy the root token printed out and set the `VAULT_TOKEN` variable.
Feel free to also set the `VAULT_ADDR` variable to allow using the normal `vault` CLI to interact with Vault as well.

```shell
$ export VAULT_TOKEN="s.dgjfnskdfgnksd"
$ export VAULT_ADDR="http://192.168.33.10:8200"
$ vault status
```

With everything set up, run the VVW tool as follows:

```shell
$ ../../bin/vvw apply -f vvw.hcl
```

If everything worked correctly then you should see a load of green "SUCCESS" messages and no red "ERROR" messages.
