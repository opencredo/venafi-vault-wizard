package download_plugin_test

import (
	"github.com/opencredo/venafi-vault-wizard/helpers/download_plugin"

	"testing"
)

func TestDownloadPluginAndUnzip(t *testing.T) {
	dl := download_plugin.NewPluginDownloader()
	_, actualSHA, err := dl.DownloadPluginAndUnzip("https://github.com/Venafi/vault-pki-backend-venafi/releases/download/v0.8.3/venafi-pki-backend_v0.8.3_linux.zip")
	if err != nil {
		t.Fatalf("Error downloading plugin: %s", err)
	}

	expectedSHA := "4440ee7d3cde5fe2aaab2f0276d645d37aef8edc86651cc183c31c22cd39ea67"

	if actualSHA != expectedSHA {
		t.Fatalf("SHAs did not match, expected (%s), got (%s)", expectedSHA, actualSHA)
	}
}
