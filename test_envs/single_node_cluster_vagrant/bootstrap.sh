#!/usr/bin/env bash
# Install, configure, and start HashiCorp Vault on multiple nodes, in HA configuration

set -o errexit
set -o pipefail
set -o nounset

__node_ip=$1

# Called at bottom of file
main() {
  install_dependencies
  prepare_vault_config
  start_vault

  echo "Everything installed and Vault service started!"
  export VAULT_ADDR=http://localhost:8200

  initialise_node

  print_root_token
}

install_dependencies() {
  # Add HashiCorp PPA
  curl --show-error --silent --fail --location https://apt.releases.hashicorp.com/gpg | apt-key add -
  apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release --codename --short) main"
  apt-get update

  # Install dependencies
  apt-get install --yes vault jq sshpass
}

prepare_vault_config() {
  # Create plugin_dir
  mkdir /etc/vault.d/plugins
  chown vagrant /etc/vault.d/plugins

  # Create Vault config file with setup for HA integrated storage
  cat > /etc/vault.d/vault.hcl <<EOF
listener "tcp" {
  address       = "0.0.0.0:8200"
  tls_disable   = true
}
plugin_directory = "/etc/vault.d/plugins"
storage "file" {
  path = "/opt/vault/data"
}
log_level = "debug"
api_addr = "http://${__node_ip}:8200"
EOF
}

start_vault() {
  # Enable Vault service
  systemctl enable vault.service
  systemctl start vault.service

  # Wait for Vault to start
  while ! nc -w 1 localhost 8200 </dev/null; do sleep 1; done
}

initialise_node() {
  echo "Initialising Vault"
  # Init and unseal Vault
  vault operator init \
      -key-shares=1 \
      -key-threshold=1 \
      -format=json \
      > /home/vagrant/init-keys.json

  unseal_vault
}

print_root_token() {
  echo "Root token:"
  jq --raw-output '.root_token' /home/vagrant/init-keys.json
}

unseal_vault() {
  local unseal_key
  unseal_key=$(jq --raw-output '.unseal_keys_b64 | first' /home/vagrant/init-keys.json)
  vault operator unseal "${unseal_key}"
}

main