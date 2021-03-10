package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddress string
	vaultToken   string
)

var rootCmd = &cobra.Command{
	Use:   "vvw",
	Short: "Venafi Vault Wizard",
	Long:  "VVW is a wizard to automate the installation and verification of Venafi PKI plugins for HashiCorp Vault.",
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&vaultAddress,
		"vaultAddress",
		"https://127.0.0.1:8200",
		"Vault HTTP API endpoint",
	)
	rootCmd.PersistentFlags().StringVar(
		&vaultToken,
		"vaultToken",
		"root",
		"Token used to authenticate with Vault",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
