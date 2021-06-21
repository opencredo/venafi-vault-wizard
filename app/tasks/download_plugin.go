package tasks

import (
	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
)

type DownloadPluginInput struct {
	Reporter reporter.Report
	Plugin   plugins.PluginConfig
}

// DownloadPlugin gets the plugin's download URL from its Impl.GetDownloadURL(), then downloads and unzips it, returning
// the plugin binary itself as a byte slice, and the SHA as a string
func DownloadPlugin(i *DownloadPluginInput) ([]byte, string, error) {
	pluginDownloadSection := i.Reporter.AddSection("Downloading plugin")

	downloadCheck := pluginDownloadSection.AddCheck("Downloading plugin...")

	pluginBytes, sha, err := i.Plugin.Impl.DownloadPlugin()
	if err != nil {
		downloadCheck.Errorf("Could not download plugin: %s", err)
		return nil, "", err
	}

	downloadCheck.Success("Successfully downloaded plugin")
	return pluginBytes, sha, nil
}
