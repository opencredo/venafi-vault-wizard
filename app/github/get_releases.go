package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type release struct {
	Tag         string    `json:"tag_name"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []*asset
}

type asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

// GetReleases will find the releases for Github repo "https://github.com/{repoOwnerAndName}". It will then either try
// to match the Git tag with desiredVersion, or if desiredVersion is blank, find the latest version. With a particular
// release selected, it will then loop through the assets attached to the release and find the first one matching
// assetSearchSubstring using strings.Contains. This function ignores drafts and prereleases. It will return the
// download URL for the particular asset, followed by the version it corresponds to (useful if desiredVersion was blank)
func GetReleases(repoOwnerAndName, desiredVersion, assetSearchSubstring string) (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", repoOwnerAndName)

	releases, err := downloadFromGithub(url)
	if err != nil {
		return "", "", err
	}

	var desiredRelease *release

	if desiredVersion == "" {
		desiredRelease = getLatestRelease(releases)
	} else {
		desiredRelease, err = getReleaseWithVersion(desiredVersion, releases)
		if err != nil {
			return "", "", err
		}
	}

	desiredAsset, err := getAssetMatchingSubstring(assetSearchSubstring, desiredRelease.Assets)
	if err != nil {
		return "", "", err
	}

	return desiredAsset.URL, desiredRelease.Tag, nil
}

func downloadFromGithub(url string) ([]*release, error) {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var releases []*release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return nil, err
	}

	releases = filterOutDraftAndPrerelease(releases)

	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found")
	}

	return releases, nil
}

func filterOutDraftAndPrerelease(releases []*release) []*release {
	var filtered []*release
	for _, r := range releases {
		if !r.Draft && !r.Prerelease {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func getLatestRelease(releases []*release) *release {
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].PublishedAt.After(releases[j].PublishedAt)
	})

	return releases[0]
}

func getReleaseWithVersion(version string, releases []*release) (*release, error) {
	for _, r := range releases {
		if r.Tag == version {
			return r, nil
		}
	}

	return nil, fmt.Errorf("no release found for version %s", version)
}

func getAssetMatchingSubstring(substr string, assets []*asset) (*asset, error) {
	for _, a := range assets {
		if strings.Contains(a.Name, substr) {
			return a, nil
		}
	}

	return nil, fmt.Errorf("no asset found with substring %s", substr)
}
