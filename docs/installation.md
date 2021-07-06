# Installation

While this tool is in development, you are required to build it yourself in order to use it.
To do this, simply run:

```shell
$ make build
```

Once this runs successfully, you can test it as follows:

```
$ ./bin/vvw
VVW is a wizard to automate the installation and verification of Venafi PKI plugins for HashiCorp Vault.

Usage:
  vvw [command]

Available Commands:
  generate-config Generates config file based on asking questions
  apply           Applies desired state as specified in config file
  help            Help about any command

Flags:
  -f, --configFile string   Path to config file to use to configure Venafi Vault plugin (default "vvw_config.hcl")
  -h, --help                help for vvw

Use "vvw [command] --help" for more information about a command.
```