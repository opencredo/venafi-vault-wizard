package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddress string
	vaultToken   string
	sshUser      string
	sshPassword  string
	sshPort      uint
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
	rootCmd.PersistentFlags().StringVar(
		&sshUser,
		"sshUser",
		"username",
		"Username with which to log into Vault server over SSH (must have sudo privileges)",
	)
	rootCmd.PersistentFlags().StringVar(
		&sshPassword,
		"sshPassword",
		"password",
		"Password for SSH user to log into Vault server with",
	)
	rootCmd.PersistentFlags().UintVar(
		&sshPort,
		"sshPort",
		22,
		"Port on which SSH is running on the Vault server",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
