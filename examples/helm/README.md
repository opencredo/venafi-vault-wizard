# Vault Test Environment with HA Cluster running on Integrated Storage (RAFT) in Kubernetes

This test environment spins up a 3-node Vault cluster in containers on Kubernetes.
It uses [Helm](https://helm.sh) and the official [vault-helm](https://github.com/hashicorp/vault-helm) chart to describe the desired Kubernetes resources.
The chart is then customised (via `values.yaml`) to gear the cluster more towards test/development instead of a production cluster.
One key difference, however, is that the container image is overridden to use a custom one with the Vault plugins baked in.
See the `Dockerfile` for how this is put together.

Depending on the Kubernetes cluster used, it might be easier to use some kind of ingress to access the Vault API, but by default it is assumed that `kubectl port-forward` will be used as it will work with local clusters.
If this is changed, then the `vvw.hcl` config file should be updated to reflect the new `api_address`.

## Setup

There are two example configuration files named `tpp-vvw.hcl` and `vaas-vvw.hcl`  that can be used
in conjunction with Venafi Trust Protection Platform and Venafi as a Service respectively.

To use the `tpp-vvw.hcl` file the following environment variables need to be set:

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

To use the `vaas-vvw.hcl` file the following environment variables need to be set:

- `VAULT_TOKEN`: this will be printed to stdout inside the VM so should be reflected in the Vagrant output
- `VENAFI_API_KEY`: this can be generated and is available through the Venafi web console under `API Keys` within `User Preferences`

If using Venafi as a Service, then adjust the `zone`s and `secret`s accordingly.

## Usage

First build the Dockerfile containing Vault and the plugins:

```shell
$ make build
```

To download the Helm dependencies and spin up the cluster, run:

```shell
$ helm dependency update
$ helm install vvwtestcluster .
```

Then to initialise and unseal Vault, run:

```shell
$ kubectl exec -it vvwtestcluster-vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > init-keys.json
$ kubectl exec -it vvwtestcluster-vault-0 -- vault operator unseal $(jq -r '.unseal_keys_b64 | first' init-keys.json)
$ echo "Root token: $(jq -r '.root_token' init-keys.json)"
```

When this has completed, copy the root token printed out and set the `VAULT_TOKEN` variable.

```shell
$ export VAULT_TOKEN="s.dgjfnskdfgnksd"
```

The Helm chart will create a number of Kubernetes services to reference the Vault pods.
The easiest way to temporarily access these without configuring anything extra is by running `kubectl port-forward` in another terminal.
This will allow accessing the API via `localhost:8200` while the command is running.
See the command below for an example of doing this.

```shell
$ kubectl port-forward service/vvwtestcluster-vault 8200
```

Feel free to also set the `VAULT_ADDR` variable to allow using the normal `vault` CLI to interact with Vault as well.

```shell
$ export VAULT_ADDR="http://localhost:8200"
$ vault status
```

With everything set up, run the VVW tool as follows for TPP:

```shell
$ ../../bin/vvw apply -f tpp-vvw.hcl
```

or the following command for VaaS:

```shell
$ ../../bin/vvw apply -f vaas-vvw.hcl
```

If everything worked correctly then you should see a load of green "SUCCESS" messages and no red "ERROR" messages.

## Cleanup

When you are finished with the cluster, the Kubernetes resources can be torndown by running:

```shell
$ helm uninstall vvwtestcluster
```
