package redisenterprise

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/github"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
	"github.com/opencredo/venafi-vault-wizard/app/vault/api"
)

func (c *RedisEnterpriseConfig) DownloadPlugin() ([]byte, string, error) {
	pluginURL, err := github.GetRelease(
		"RedisLabs/vault-plugin-database-redis-enterprise",
		c.Version,
		"linux_amd64",
	)
	if err != nil {
		return nil, "", err
	}

	pluginBytes, err := downloader.DownloadFile(pluginURL)
	if err != nil {
		return nil, "", err
	}

	return pluginBytes, downloader.GetSHAString(pluginBytes), nil
}

func (c *RedisEnterpriseConfig) Configure(report reporter.Report, vaultClient api.VaultAPIClient) error {
	configurePluginSection := report.AddSection("Setting up redisenterprise")

	for _, role := range c.Roles {
		dbCheck := configurePluginSection.AddCheck("Adding Redis database...")

		dbPath := fmt.Sprintf("%s/config/%s", c.MountPath, role.DBCluster.Name)
		_, err := vaultClient.WriteValue(
			dbPath,
			map[string]interface{}{
				"plugin_name":   c.PluginName,
				"url":           role.DBCluster.URL,
				"username":      role.DBCluster.Username,
				"password":      role.DBCluster.Password,
				"database":      role.DBCluster.DatabaseName,
				"allowed_roles": role.Name,
			},
		)
		if err != nil {
			dbCheck.Errorf("Error adding Redis database: %s", err)
			return err
		}

		dbCheck.Successf("Configured Redis database at %s", dbPath)

		roleCheck := configurePluginSection.AddCheck("Adding Redis role...")

		rolePath := fmt.Sprintf("%s/roles/%s", c.MountPath, role.Name)
		creationStatements := fmt.Sprintf(`{"role": "%s"}`, role.DBRole)
		_, err = vaultClient.WriteValue(
			rolePath,
			map[string]interface{}{
				"db_name":             role.DBCluster.Name,
				"creation_statements": creationStatements,
			},
		)
		if err != nil {
			roleCheck.Errorf("Error adding Redis role: %s", err)
			return err
		}

		roleCheck.Successf("Configured Redis role at %s", rolePath)
	}
	return nil
}

func (c *RedisEnterpriseConfig) Check(report reporter.Report, vaultClient api.VaultAPIClient) error {
	for _, role := range c.Roles {
		roleIssuePath := fmt.Sprintf("%s/creds/%s", c.MountPath, role.Name)

		section := report.AddSection("Testing Redis role " + role.Name)

		check := section.AddCheck(fmt.Sprintf("Requesting test credentials from %s", roleIssuePath))

		creds, err := vaultClient.ReadValue(roleIssuePath)
		if err != nil {
			check.Errorf("Error requesting credentials from %s: %s", roleIssuePath, err)
			return err
		}

		check.Successf(
			"Successfully retrieved test credentials from Redis role %s, with username %s",
			role.Name,
			creds["username"],
		)
	}

	return nil
}
