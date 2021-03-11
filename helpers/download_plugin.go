package helpers

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DownloadPluginAndUnzip downloads the ZIP archive from the given URL, unzips it, and returns the plugin and its SHA
func DownloadPluginAndUnzip(url string) (plugin []byte, sha string, err error) {
	pluginBytes, err := downloadZipFile(url)
	if err != nil {
		return nil, "", err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(pluginBytes), int64(len(pluginBytes)))
	if err != nil {
		return nil, "", err
	}

	// Check we've got two files, one for the plugin binary, and one for the SHA
	if len(zipReader.File) != 2 {
		return nil, "", fmt.Errorf("expected 2 files in the plugin's zip file, got %d", len(zipReader.File))
	}

	var expectedSHA string

	// Read SHA file to string and return plugin binary as reader
	for _, file := range zipReader.File {
		unzippedBytes, err := readZippedFile(file)
		if err != nil {
			return nil, "", err
		}

		if strings.Contains(file.Name, "SHA256SUM") {
			expectedSHA = strings.TrimSpace(string(unzippedBytes))
		} else {
			plugin = unzippedBytes
		}
	}

	actualSHA := getSHAString(plugin)

	err = checkSHAsMatch(expectedSHA, actualSHA)
	if err != nil {
		return nil, "", err
	}

	return plugin, expectedSHA, nil
}

func downloadZipFile(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	pluginBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("uh oh")
	}

	return pluginBytes, nil
}

func readZippedFile(file *zip.File) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	unzipped, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return unzipped, nil
}

func getSHAString(file []byte) string {
	rawHash := sha256.Sum256(file)
	return hex.EncodeToString(rawHash[:])
}

func checkSHAsMatch(expected, actual string) error {
	if strings.Compare(expected, actual) != 0 {
		return fmt.Errorf("expected SHA checksums to match")
	}

	return nil
}
