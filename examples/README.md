# VVW Example Environments

This directory contains a number of test environments that can be used to spin up Vault clusters in various configurations.
Each one should contain another `README` explaining more.
They should also contain a `vvw.hcl` configuration file that works with the environment to allow testing VVW in each scenario, but also to be adapted to similar scenarios.
For example, if someone wants to use the Venafi plugins with Vault running in Kubernetes, the `vvw.hcl` file in that test environment should be a good starting point.

The currently supported test environments are:

- `single_node_cluster_vagrant` - Single node Vault server running in a VM
- `integrated_storage_ha_cluster_vagrant` - 3-node Vault cluster using integrated Raft storage mode running in VMs
- `helm` - 3-node Vault cluster using integrated Raft storage mode running in containers on Kubernetes
- `consul_ha_cluster_vagrant` - similar to `integrated_storage_ha_cluster_vagrant` but includes a Consul server in another VM and uses that for HA storage