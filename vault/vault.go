package vault

import (
	"fmt"
	vaultAPI "github.com/hashicorp/vault/api"
)

// Vault represents a HashiCorp Vault instance and the operations available on it
type Vault interface {
	// GetPluginDir queries the server for the local plugin directory
	GetPluginDir() (directory string, err error)
	// CopyFile connects the Vault server via SSH and writes some random text to a file (will eventually copy the plugin)
	CopyFile(destination string) error
}

type vault struct {
	Config *Config
	Client *vaultAPI.Client
}

// Config represents the configuration values needed to connect to Vault via the API and SSH
type Config struct {
	// Address of the Vault server that the API is served on. Equivalent of setting VAULT_ADDR for the vault CLI
	APIAddress string
	// Authentication token to perform Vault operations. Must have sufficient permissions
	Token string
	// Address of the Vault server that can be used for SSH access.
	SSHAddress string
}

// NewVault returns an instance of the Vault client
func NewVault(config *Config) (Vault, error) {
	client, err := vaultAPI.NewClient(&vaultAPI.Config{
		Address: config.APIAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting vault client: %w", ErrInvalidAddress)
	}

	client.SetToken(config.Token)
	return &vault{config, client}, nil
}

func (v *vault) GetPluginDir() (string, error) {
	secret, err := v.Client.Logical().Read("sys/config/state/sanitized")
	if err != nil {
		// TODO: parse out error codes and adjust error message accordingly
		return "", err
	}

	dir, ok := secret.Data["plugin_directory"].(string)
	if !ok {
		return "", ErrPluginDirNotConfigured
	}

	return dir, nil
}
