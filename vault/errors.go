package vault

import "errors"

var ErrInvalidAddress = errors.New("can't access Vault at address provided")
var ErrPluginDirNotConfigured = errors.New("plugin_directory not set in Vault config file")
var ErrReadingVaultPath = errors.New("cannot read Vault path")
