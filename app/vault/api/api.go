package api

import (
	"fmt"

	vaultAPI "github.com/hashicorp/vault/api"
	vaultConsts "github.com/hashicorp/vault/sdk/helper/consts"

	"github.com/opencredo/venafi-vault-wizard/app/vault"
	"github.com/opencredo/venafi-vault-wizard/app/vault/lib"
)

// VaultAPIClient represents a HashiCorp Vault instance and the operations available on it via the Vault API. For
// operations involving SSH, see the vault/ssh/VaultSSHClient interface instead.
type VaultAPIClient interface {
	// GetPluginDir queries the server for the local plugin directory
	GetPluginDir() (directory string, err error)
	// RegisterPlugin adds the plugin to the VaultPlugin Catalog
	RegisterPlugin(name, command, sha string) error
	// GetPlugin returns information about a registered plugin (command, sha, args etc)
	GetPlugin(name string) (map[string]interface{}, error)
	// ReloadPlugin reloads a plugin (globally across a cluster if Vault is clustered) and waits for the number of
	// completed reloads to equal the number of replicas
	ReloadPlugin(name string) error
	// MountPlugin mounts a secret engine at the specified path. Equivalent to vault secrets enable -plugin-name=name -path=path
	MountPlugin(name, path string) error
	// GetMountPluginName checks which backend is used for particular mount
	GetMountPluginName(path string) (string, error)
	// WriteValue writes to the specified path. Equivalent to `$ vault write path value1=v1 value2=v2`
	WriteValue(path string, value map[string]interface{}) (map[string]interface{}, error)
	// ReadValue reads from the specified path. Equivalent to `$ vault read path`
	ReadValue(path string) (map[string]interface{}, error)
	// GetVaultConfig reads the config from sys/config/state/sanitized and returns it as a map
	GetVaultConfig() (map[string]interface{}, error)
	// IsMLockDisabled checks to see if the server was run with the disable_mlock option
	IsMLockDisabled() (bool, error)
}

type vaultAPIClient struct {
	Config      *Config
	VaultClient lib.VaultAPIWrapper
}

// Config represents the configuration values needed to connect to Vault via the API
type Config struct {
	// Address of the Vault server that the API is served on. Equivalent of setting VAULT_ADDR for the vault CLI
	APIAddress string
	// Authentication token to perform Vault operations. Must have sufficient permissions
	Token string
}

// NewClient returns an instance of the Vault API client
func NewClient(config *Config, apiClient lib.VaultAPIWrapper) (VaultAPIClient, error) {
	err := apiClient.SetAddress(config.APIAddress)
	if err != nil {
		return nil, err
	}
	apiClient.SetToken(config.Token)

	return &vaultAPIClient{config, apiClient}, nil
}

func (v *vaultAPIClient) GetPluginDir() (string, error) {
	config, err := v.GetVaultConfig()
	if err != nil {
		return "", err
	}

	dir, ok := config["plugin_directory"].(string)
	if dir == "" || !ok {
		return "", vault.ErrPluginDirNotConfigured
	}

	return dir, nil
}

func (v *vaultAPIClient) RegisterPlugin(name, command, sha string) error {
	err := v.VaultClient.RegisterPlugin(&vaultAPI.RegisterPluginInput{
		Name:    name,
		Type:    vaultConsts.PluginTypeSecrets,
		Command: command,
		SHA256:  sha,
	})
	if err != nil {
		return fmt.Errorf("error writing sys/plugins/catalog/secret: %w", err)
	}

	return nil
}

func (v *vaultAPIClient) GetPlugin(name string) (map[string]interface{}, error) {
	plugin, err := v.VaultClient.GetPlugin(&vaultAPI.GetPluginInput{
		Name: name,
		Type: vaultConsts.PluginTypeSecrets,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting details for plugin %s: %w", name, err)
	}

	return map[string]interface{}{
		"command": plugin.Command,
		"args":    plugin.Args,
		"sha":     plugin.SHA256,
	}, nil
}

func (v *vaultAPIClient) ReloadPlugin(name string) error {
	reloadID, err := v.VaultClient.ReloadPlugin(&vaultAPI.ReloadPluginInput{
		Plugin: name,
	})
	if err != nil {
		return fmt.Errorf("error reloading plugin %s: %w", name, err)
	}

	if reloadID == "" {
		return nil
	}

	return fmt.Errorf("wasn't expecting a reload ID from the plugin reload, are there multiple clusters?")
}

func (v *vaultAPIClient) MountPlugin(name, path string) error {
	err := v.VaultClient.Mount(path, &vaultAPI.MountInput{
		Type: name,
	})
	if err != nil {
		// TODO: check for "Unrecognized remote plugin message" and see whether it's mlock or api_addr
		return fmt.Errorf("error mounting plugin %s at path %s: %w", name, path, err)
	}

	return nil
}

func (v *vaultAPIClient) GetMountPluginName(path string) (string, error) {
	mounts, err := v.VaultClient.ListMounts()
	if err != nil {
		return "", fmt.Errorf("error listing mounts: %w", err)
	}

	mount, ok := mounts[path+"/"]
	if !ok {
		return "", fmt.Errorf("nothing mounted at path %s: %w", path, vault.ErrPluginNotMounted)
	}

	return mount.Type, nil
}

func (v *vaultAPIClient) WriteValue(path string, value map[string]interface{}) (map[string]interface{}, error) {
	secret, err := v.VaultClient.Write(path, value)
	if err != nil {
		return nil, fmt.Errorf("error writing to path %s: %w", path, err)
	}

	return secret, nil
}

func (v *vaultAPIClient) ReadValue(path string) (map[string]interface{}, error) {
	secret, err := v.VaultClient.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error reading from path %s: %w", path, err)
	}

	return secret, nil
}

func (v *vaultAPIClient) GetVaultConfig() (map[string]interface{}, error) {
	return v.ReadValue("sys/config/state/sanitized")
}

func (v *vaultAPIClient) IsMLockDisabled() (bool, error) {
	config, err := v.GetVaultConfig()
	if err != nil {
		return false, err
	}

	disabled, ok := config["disable_mlock"].(bool)
	if !ok {
		return false, fmt.Errorf("error, `disable_mlock` option not found in config")
	}

	return disabled, nil
}
