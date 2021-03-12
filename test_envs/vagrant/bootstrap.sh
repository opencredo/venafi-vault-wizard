#!/usr/bin/env bash

# Install Vault
curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -
apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
apt-get update
apt-get install -y vault

# Set up filesystem for Vault
mkdir /etc/vault.d/plugins
chown vagrant /etc/vault.d/plugins
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
api_addr = "http://0.0.0.0:8200"
EOF

# Enable Vault service
systemctl enable vault.service
systemctl start vault.service

apt-get install -y jq
export VAULT_ADDR=http://localhost:8200

# Init and unseal Vault
vault operator init \
    -key-shares=1 \
    -key-threshold=1 \
    -format=json \
    > /home/vagrant/init-keys.json

UNSEAL_KEY=$(jq -r '.unseal_keys_b64 | first' /home/vagrant/init-keys.json)
vault operator unseal ${UNSEAL_KEY}

jq -r '.root_token' /home/vagrant/init-keys.json
