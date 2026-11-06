package internal

import (
	"context"
	"time"
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// UpdateInfo describes the result of an update check.
type UpdateInfo struct {
	Current    string
	Latest     string
	Available  bool
	ReleaseURL string
}

// isUpdateAvailable reports whether latest is strictly newer than current.
func isUpdateAvailable(current, latest string) bool {
	cur, err := semver.Make(current)
	if err != nil {
		// Unknown current version (e.g. "dev") counts as updatable.
		return true
	}
	lat, err := semver.Make(latest)
	if err != nil {
		return false
	}
	return lat.GT(cur)
}

// CheckForUpdate queries GitHub releases for the latest version.
func CheckForUpdate(ctx context.Context) (UpdateInfo, error) {
	info := UpdateInfo{Current: Version}
	slug := RepoOwner + "/" + RepoName

	release, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return info, err
	}
	if !found {
		return info, nil
	}

	info.Latest = release.Version.String()
	info.ReleaseURL = release.URL

	info.Available = isUpdateAvailable(Version, release.Version.String())
	return info, nil
}

// PerformUpdate downloads and installs the latest release, replacing the
// running binary. It returns the resulting UpdateInfo.
func PerformUpdate(ctx context.Context) (UpdateInfo, error) {
	info, err := CheckForUpdate(ctx)
	if err != nil {
		return info, err
	}
	if !info.Available {
		return info, nil
	}

	current, _ := semver.Make(Version)
	slug := RepoOwner + "/" + RepoName
	release, err := selfupdate.UpdateSelf(current, slug)
	if err != nil {
		return info, fmt.Errorf("update failed: %w", err)
	}
	info.Latest = release.Version.String()
	info.Available = false
	return info, nil
}

// DefaultUpdateTimeout is used for best-effort update checks.
const DefaultUpdateTimeout = 5 * time.Second
