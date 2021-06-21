package downloader

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

// DownloadPluginAndUnzip is a helper method that uses the other methods in this package. It takes a URL referring to a
// Vault plugin ZIP file, downloads it with DownloadFile, unzips it with UnzipPluginAndSHA, verifies the SHA is as
// expected using CheckSHAsMatch, and then returns the resulting byte slice with the plugin and string with its SHA
func DownloadPluginAndUnzip(url string) ([]byte, string, error) {
	pluginBytes, err := DownloadFile(url)
	if err != nil {
		return nil, "", err
	}

	plugin, expectedSHA, err := UnzipPluginAndSHA(pluginBytes)
	if err != nil {
		return nil, "", err
	}

	err = CheckSHAsMatch(plugin, expectedSHA)
	if err != nil {
		return nil, "", err
	}

	return plugin, expectedSHA, nil
}

func DownloadFile(url string) ([]byte, error) {
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

// UnzipPluginAndSHA reads a byte slice containing a zipped archive and returns the plugin binary file as a byte slice,
// and the plugin SHA as a string. It expects two files inside the archive: one with "SHA256SUM" in its name, and the
// other containing the plugin itself.
func UnzipPluginAndSHA(zipFile []byte) ([]byte, string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipFile), int64(len(zipFile)))
	if err != nil {
		return nil, "", err
	}

	// Check we've got two files, one for the plugin binary, and one for the SHA
	if len(zipReader.File) != 2 {
		return nil, "", fmt.Errorf("expected 2 files in the plugin's zip file, got %d", len(zipReader.File))
	}

	var plugin []byte
	var expectedSHA string

	// Read SHA file to string and get plugin bytes
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

	return plugin, expectedSHA, nil
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

// CheckSHAsMatch takes a byte slice and an expected SHA string, and verifies that the checksum matches the byte slice
func CheckSHAsMatch(plugin []byte, expected string) error {
	actualSHA := GetSHAString(plugin)
	if strings.Compare(expected, actualSHA) != 0 {
		return fmt.Errorf("expected SHA checksums to match")
	}

	return nil
}

// GetSHAString takes a file as a byte slice and returns its SHA256 checksum
func GetSHAString(file []byte) string {
	rawHash := sha256.Sum256(file)
	return hex.EncodeToString(rawHash[:])
}
