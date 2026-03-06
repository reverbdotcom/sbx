package update

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/go-github/v67/github"
	gh "github.com/reverbdotcom/sbx/github"
	"github.com/reverbdotcom/sbx/version"
)

const owner = "reverbdotcom"
const repo = "sbx"

var getLatestReleaseFn = _getLatestRelease
var executableFn = os.Executable
var downloadAndReplaceFn = _downloadAndReplace

func _getLatestRelease() (*github.RepositoryRelease, error) {
	client, err := gh.Client()
	if err != nil {
		return nil, err
	}

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func Run() (string, error) {
	current := version.Get()

	release, err := getLatestReleaseFn()
	if err != nil {
		return "", fmt.Errorf("failed to check for updates: %v", err)
	}

	latest := release.GetTagName()

	if latest == "" || latest == current {
		return "", nil
	}

	assetName := fmt.Sprintf("sbx-darwin-%s.tar.gz", runtime.GOARCH)

	var assetURL string
	for _, asset := range release.Assets {
		if asset.GetName() == assetName {
			assetURL = asset.GetBrowserDownloadURL()
			break
		}
	}

	if assetURL == "" {
		return "", fmt.Errorf("failed to check for updates: no asset %s in release %s", assetName, latest)
	}

	fmt.Printf("updating sbx %s -> %s...\n", current, latest)

	err = downloadAndReplaceFn(assetURL)
	if err != nil {
		return "", fmt.Errorf("failed to update: %v", err)
	}

	return "updated!\n", nil
}

func _downloadAndReplace(assetURL string) error {
	client, err := gh.Client()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", assetURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/octet-stream")

	httpClient := client.Client()
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d downloading %s", resp.StatusCode, assetURL)
	}

	binary, err := extractBinary(resp.Body)
	if err != nil {
		return err
	}

	return replaceBinary(binary)
}

func extractBinary(r io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("gzip: %v", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tar: %v", err)
		}

		if hdr.Name == "sbx" && hdr.Typeflag == tar.TypeReg {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("reading binary from tar: %v", err)
			}
			return data, nil
		}
	}

	return nil, fmt.Errorf("sbx binary not found in tarball")
}

func replaceBinary(newBinary []byte) error {
	exe, err := executableFn()
	if err != nil {
		return fmt.Errorf("finding executable path: %v", err)
	}

	resolved, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return fmt.Errorf("resolving symlinks: %v", err)
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return fmt.Errorf("stat current binary: %v", err)
	}

	dir := filepath.Dir(resolved)
	tmp, err := os.CreateTemp(dir, "sbx-update-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %v", err)
	}
	tmpPath := tmp.Name()

	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()

	if _, err := tmp.Write(newBinary); err != nil {
		return fmt.Errorf("writing temp file: %v", err)
	}

	if err := tmp.Chmod(info.Mode()); err != nil {
		return fmt.Errorf("setting permissions: %v", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %v", err)
	}

	if err := os.Rename(tmpPath, resolved); err != nil {
		return fmt.Errorf("replacing binary: %v", err)
	}

	return nil
}
