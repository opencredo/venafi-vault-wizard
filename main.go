package main

import (
	"fmt"
	"log"

	"github.com/opencredo/venafi-vault-wizard/vault"
)

func main() {
	vault, err := vault.NewVault(&vault.Config{
		APIAddress: "http://localhost:8200",
		Token:      "",
		SSHAddress: "localhost:2222",
	})
	if err != nil {
		log.Fatalf("Error getting vault instance %s", err)
	}

	dir, err := vault.GetPluginDir()
	if err != nil {
		log.Fatalf("Error getting plugin dir %s", err)
	}

	err = vault.CopyFile("/config/randomfile.txt")
	if err != nil {
		log.Fatalf("Error copying file %s", err)
	}

	fmt.Printf("Hello Venafi Cloud and HashiCorp Vault, plugin dir is %s\n", dir)
}
