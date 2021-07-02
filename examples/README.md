# VVW Example Environments

This directory contains a number of test environments that can be used to spin up Vault clusters in various configurations.
Each example contains a `README` explaining more about how it is built and executed.  Alongside the `README` each example also 
contains two configuration files named `tpp-vvw.hcl` and `vaas-vvw.hcl`.  Each configuration allows VVW to be tested with either the Venafi Trust Protection Platform or Venafi as a Service.
The aim is to provide HCL configuration files that can be used as a guide and adapted to similar scenarios.
For example, if someone wants to use the Venafi plugins with Vault running in Kubernetes, the HCL files in that example should be a good starting point.

The currently supported test environments are:

- `single_node_cluster_vagrant` - Single node Vault server running in a VM
- `integrated_storage_ha_cluster_vagrant` - 3-node Vault cluster using integrated Raft storage mode running in VMs
- `helm` - 3-node Vault cluster using integrated Raft storage mode running in containers on Kubernetes
- `consul_ha_cluster_vagrant` - similar to `integrated_storage_ha_cluster_vagrant` but includes a Consul server in another VM and uses that for HA storage