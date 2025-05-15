package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

var Version = "dev"

func GetVersion() string {
	// First priority: version set by ldflags (GoReleaser)
	if Version != "dev" {
		return Version
	}

	// Second priority: Get from build info
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	// Third priority: Check if installed via go install with version
	if bi.Main.Version != "(devel)" && bi.Main.Version != "" {
		return bi.Main.Version
	}

	// Fourth priority: Try to get from VCS info
	var vcsRevision string
	var vcsTime time.Time

	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.time":
			vcsTime, _ = time.Parse(time.RFC3339, setting.Value)
		case "vcs.tag":
			if setting.Value != "" {
				return setting.Value // Return tag if available
			}
		}
	}

	// Return commit info if available
	if vcsRevision != "" {
		return fmt.Sprintf("%s (%s)", vcsRevision[:8], vcsTime.Format("2006-01-02"))
	}

	// Default fallback
	return "dev"
}

func CheckForNewerVersion(currentVersion string) (bool, string, error) {
	const repoURL = "https://api.github.com/repos/mcphee11/mcphee11-tui/releases/latest"

	// Make HTTP request to GitHub API
	resp, err := http.Get(repoURL)
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse JSON response
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", fmt.Errorf("failed to parse release data: %w", err)
	}

	// Compare versions
	if release.TagName > currentVersion {
		return true, release.TagName, nil
	}

	return false, "", nil
}
