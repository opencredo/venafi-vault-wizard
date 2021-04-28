package tasks

import (
	"fmt"

	"github.com/opencredo/venafi-vault-wizard/app/downloader"
	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/reporter"
)

type DownloadPluginInput struct {
	Downloader downloader.PluginDownloader
	Reporter   reporter.Report
	Plugin     plugins.Plugin
}

// DownloadPlugin gets the plugin's download URL from its Impl.GetDownloadURL(), then downloads and unzips it, returning
// the plugin binary itself as a byte slice, and the SHA as a string
func DownloadPlugin(i *DownloadPluginInput) ([]byte, string, error) {
	pluginDownloadSection := i.Reporter.AddSection("Downloading plugin")

	downloadCheck := pluginDownloadSection.AddCheck("Downloading plugin...")

	pluginURL, err := i.Plugin.Impl.GetDownloadURL()
	if err != nil {
		downloadCheck.Error(fmt.Sprintf("Error getting plugin download URL: %s", err))
		return nil, "", err
	}

	pluginBytes, sha, err := i.Downloader.DownloadPluginAndUnzip(pluginURL)
	if err != nil {
		downloadCheck.Error(fmt.Sprintf("Could not download plugin from %s: %s", pluginURL, err))
		return nil, "", err
	}

	downloadCheck.Success("Successfully downloaded plugin")
	return pluginBytes, sha, nil
}
