package vault

import "errors"

var (
	ErrUnauthorised           = errors.New("the Vault request was unauthorised, check whether the token is valid")
	ErrNotFound               = errors.New("path not found, either it was incorrect or a backend isn't mounted properly")
	ErrVaultSealed            = errors.New("the Vault server is either sealed or down for maintenance")
	ErrPluginDirNotConfigured = errors.New("plugin_directory not set in Vault config file")
	ErrInvalidAddress         = errors.New("can't access Vault at the address provided, either the address/port is incorrect or the host is unreachable")
	ErrTLSDisabled            = errors.New("attempted to access Vault using TLS but it returned an HTTP response, check whether the protocol in the Vault address matches what's configured")
	ErrMountPathInUse         = errors.New("the mount path is in already in use")
)
