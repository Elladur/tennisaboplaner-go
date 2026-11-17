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
	Updated    bool
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
func CheckForUpdate(_ context.Context) (UpdateInfo, error) {
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
	if err := os.Chmod(tmpName, 0755); err != nil { //nolint:gosec
		_ = os.Remove(tmpName)
		return info, err
	}

	// Replace current binary.
	if err := os.Rename(tmpName, cmdPath); err != nil {
		_ = os.Remove(tmpName)
		return info, fmt.Errorf("replace failed: %w", err)
	}

	info.Latest = release.Version.String()
	info.Current = release.Version.String()
	info.Available = false
	info.Updated = true
	return info, nil
}

// isSafeArchivePath reports whether an archive entry name does not escape
// the destination directory.
func isSafeArchivePath(name string) bool {
	if filepath.IsAbs(name) {
		return false
	}
	cleaned := filepath.Clean(name)
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return false
	}
	return true
}

// isBinaryCandidate reports whether name looks like the shipped executable
// rather than an auxiliary file (LICENSE, README, settings.json, ...).
func isBinaryCandidate(name string) bool {
	base := strings.ToLower(filepath.Base(name))
	switch base {
	case "license", "readme", "readme.md", "settings.json":
		return false
	}
	switch filepath.Ext(base) {
	case ".md", ".txt", ".json", ".yml", ".yaml", ".toml":
		return false
	}
	return true
}

// extractExecutable writes the executable file found in the archive
// (or the raw body if not an archive) to dst.
//
//nolint:gocyclo
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
		// Collect candidate binary entries.
		type entry struct {
			f    *zip.File
			size uint64
		}
		var cands []entry
		for _, f := range z.File {
			if f.FileInfo().IsDir() || !isBinaryCandidate(f.Name) || !isSafeArchivePath(f.Name) {
				continue
			}
			cands = append(cands, entry{f: f, size: f.UncompressedSize64})
		}
		if len(cands) == 0 {
			return fmt.Errorf("no executable in zip")
		}
		// Prefer name match, then largest.
		chosen := cands[0]
		for _, c := range cands {
			base := strings.TrimSuffix(filepath.Base(c.f.Name), filepath.Ext(c.f.Name))
			if strings.EqualFold(base, "tennisaboplaner-go") {
				chosen = c
				break
			}
		}
		if !strings.EqualFold(strings.TrimSuffix(filepath.Base(chosen.f.Name), filepath.Ext(chosen.f.Name)), "tennisaboplaner-go") {
			// fall back to largest
			for _, c := range cands {
				if c.size > chosen.size {
					chosen = c
				}
			}
		}
		rc, err := chosen.f.Open()
		if err != nil {
			return err
		}
		defer func() { _ = rc.Close() }()
		out, err := os.Create(dst) //nolint:gosec
		if err != nil {
			return err
		}
		_, err = io.Copy(out, rc) //nolint:gosec
		_ = out.Close()
		return err
	case strings.HasSuffix(lower, ".tar.gz"), strings.HasSuffix(lower, ".tgz"):
		gz, err := gzip.NewReader(src)
		if err != nil {
			return err
		}
		defer func() { _ = gz.Close() }()
		tr := tar.NewReader(gz)
		type entry struct {
			header *tar.Header
			data   []byte
		}
		var cands []entry
		for {
			h, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if h.Typeflag != tar.TypeReg || !isBinaryCandidate(h.Name) || !isSafeArchivePath(h.Name) {
				continue
			}
			buf, err := io.ReadAll(tr)
			if err != nil {
				return err
			}
			cands = append(cands, entry{header: h, data: buf})
		}
		if len(cands) == 0 {
			return fmt.Errorf("no regular file in tar")
		}
		chosen := cands[0]
		for _, c := range cands {
			base := strings.TrimSuffix(filepath.Base(c.header.Name), filepath.Ext(c.header.Name))
			if strings.EqualFold(base, "tennisaboplaner-go") {
				chosen = c
				break
			}
		}
		if !strings.EqualFold(strings.TrimSuffix(filepath.Base(chosen.header.Name), filepath.Ext(chosen.header.Name)), "tennisaboplaner-go") {
			for _, c := range cands {
				if c.header.Size > chosen.header.Size {
					chosen = c
				}
			}
		}
		out, err := os.Create(dst) //nolint:gosec
		if err != nil {
			return err
		}
		_, err = out.Write(chosen.data)
		_ = out.Close()
		return err
	default:
		out, err := os.Create(dst) //nolint:gosec
		if err != nil {
			return err
		}
		defer func() { _ = out.Close() }()
		_, err = io.Copy(out, src) //nolint:gosec
		return err
	}
}

// DefaultUpdateTimeout is used for best-effort update checks.
const DefaultUpdateTimeout = 5 * time.Second
