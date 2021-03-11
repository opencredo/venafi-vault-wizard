package vault

import (
	"fmt"
	vaultAPI "github.com/hashicorp/vault/api"
	vaultConsts "github.com/hashicorp/vault/sdk/helper/consts"
	"io"
)

// Vault represents a HashiCorp Vault instance and the operations available on it
type Vault interface {
	// GetPluginDir queries the server for the local plugin directory
	GetPluginDir() (directory string, err error)
	// WriteFile connects the Vault server via SSH and writes some random text to a file (will eventually copy the plugin)
	WriteFile(sourceFile io.Reader, hostDestination string) error
	// RegisterPlugin adds the plugin to the Vault Plugin Catalog
	RegisterPlugin(name, command, sha string) error
	// MountPlugin mounts a secret engine at the specified path. Equivalent to vault secrets enable -plugin-name=name -path=path
	MountPlugin(name, path string) error
	// WriteValue writes to the specified path. Equivalent to `$ vault write path value1=v1 value2=v2`
	WriteValue(path string, value map[string]interface{}) error
	// ReadValue reads from the specified path. Equivalent to `$ vault read path`
	ReadValue(path string) (map[string]interface{}, error)
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
		return "", fmt.Errorf("error reading sys/config/state: %w", ErrReadingVaultPath)
	}

	dir, ok := secret.Data["plugin_directory"].(string)
	if !ok {
		return "", ErrPluginDirNotConfigured
	}

	return dir, nil
}

func (v *vault) RegisterPlugin(name, command, sha string) error {
	err := v.Client.Sys().RegisterPlugin(&vaultAPI.RegisterPluginInput{
		Name:    name,
		Type:    vaultConsts.PluginTypeSecrets,
		Command: command,
		SHA256:  sha,
	})
	if err != nil {
		// TODO: parse out error codes and adjust error message accordingly
		return fmt.Errorf("error writing sys/plugins/catalog/secret: %w", ErrWritingVaultPath)
	}

	return nil
}

func (v *vault) MountPlugin(name, path string) error {
	err := v.Client.Sys().Mount(path, &vaultAPI.MountInput{
		Type: name,
	})
	if err != nil {
		// TODO: parse out error codes and adjust error message accordingly
		return fmt.Errorf("error mounting plugin %s at path %s: %w", name, path, err)
	}

	return nil
}

func (v *vault) WriteValue(path string, value map[string]interface{}) error {
	_, err := v.Client.Logical().Write(path, value)
	if err != nil {
		// TODO: parse out error codes and adjust error message accordingly
		return err
	}

	return nil
}

func (v *vault) ReadValue(path string) (map[string]interface{}, error) {
	secret, err := v.Client.Logical().Read(path)
	if err != nil {
		// TODO: parse out error codes and adjust error message accordingly
		return nil, err
	}

	return secret.Data, nil
}
