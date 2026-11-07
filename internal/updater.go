package internal

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	slug := RepoOwner + "/" + RepoName
	release, found, err := selfupdate.DetectLatest(slug)
	if err != nil {
		return info, err
	}
	if !found {
		return info, fmt.Errorf("no release found")
	}

	// Download the release asset.
	resp, err := http.Get(release.AssetURL)
	if err != nil {
		return info, fmt.Errorf("download failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return info, fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	// Extract the executable from the archive into a temp file next to the
	// current binary so the final rename is atomic.
	cmdPath, err := os.Executable()
	if err != nil {
		return info, err
	}
	dir := filepath.Dir(cmdPath)
	tmp, err := os.CreateTemp(dir, "tennisaboplaner-go-update-*")
	if err != nil {
		return info, err
	}
	tmpName := tmp.Name()
	_ = tmp.Close()

	if err := extractExecutable(resp.Body, release.AssetURL, tmpName); err != nil {
		_ = os.Remove(tmpName)
		return info, err
	}
	if err := os.Chmod(tmpName, 0755); err != nil {
		_ = os.Remove(tmpName)
		return info, err
	}

	// Replace current binary.
	if err := os.Rename(tmpName, cmdPath); err != nil {
		_ = os.Remove(tmpName)
		return info, fmt.Errorf("replace failed: %w", err)
	}

	info.Latest = release.Version.String()
	info.Available = false
	return info, nil
}

// extractExecutable writes the first executable file found in the archive
// (or the raw body if not an archive) to dst.
func extractExecutable(src io.Reader, url, dst string) error {
	lower := strings.ToLower(url)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		buf, err := io.ReadAll(src)
		if err != nil {
			return err
		}
		r := bytes.NewReader(buf)
		z, err := zip.NewReader(r, int64(len(buf)))
		if err != nil {
			return err
		}
		for _, f := range z.File {
			if f.FileInfo().IsDir() {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer func() { _ = rc.Close() }()
			out, err := os.Create(dst)
			if err != nil {
				return err
			}
			_, err = io.Copy(out, rc)
			_ = out.Close()
			return err
		}
		return fmt.Errorf("no file in zip")
	case strings.HasSuffix(lower, ".tar.gz"), strings.HasSuffix(lower, ".tgz"):
		gz, err := gzip.NewReader(src)
		if err != nil {
			return err
		}
		defer func() { _ = gz.Close() }()
		tr := tar.NewReader(gz)
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if h.Typeflag != tar.TypeReg {
				continue
			}
			out, err := os.Create(dst)
			if err != nil {
				return err
			}
			_, err = io.Copy(out, tr)
			_ = out.Close()
			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("no regular file in tar")
	default:
		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer func() { _ = out.Close() }()
		_, err = io.Copy(out, src)
		return err
	}
}

// DefaultUpdateTimeout is used for best-effort update checks.
const DefaultUpdateTimeout = 5 * time.Second
