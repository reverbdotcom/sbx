package update

import (
	"errors"
	"testing"

	"github.com/google/go-github/v67/github"
)

func makeRelease(tag string, assets []*github.ReleaseAsset) *github.RepositoryRelease {
	return &github.RepositoryRelease{
		TagName: &tag,
		Assets:  assets,
	}
}

func makeAsset(name, url string) *github.ReleaseAsset {
	return &github.ReleaseAsset{
		Name:               &name,
		BrowserDownloadURL: &url,
	}
}

func TestRun(t *testing.T) {
	t.Run("returns empty when already up to date", func(t *testing.T) {
		getLatestReleaseFn = func() (*github.RepositoryRelease, error) {
			tag := "v1.22.0" // same as SBX_VERSION
			return makeRelease(tag, nil), nil
		}

		out, err := Run()

		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if out != "" {
			t.Errorf("got %q, want empty", out)
		}
	})

	t.Run("returns error when release check fails", func(t *testing.T) {
		getLatestReleaseFn = func() (*github.RepositoryRelease, error) {
			return nil, errors.New("network error")
		}

		_, err := Run()

		if err == nil {
			t.Error("got nil, want error")
		}
	})

	t.Run("returns error when no matching asset", func(t *testing.T) {
		getLatestReleaseFn = func() (*github.RepositoryRelease, error) {
			return makeRelease("v99.0.0", nil), nil
		}

		_, err := Run()

		if err == nil {
			t.Error("got nil, want error")
		}
	})

	t.Run("returns empty when tag is empty", func(t *testing.T) {
		getLatestReleaseFn = func() (*github.RepositoryRelease, error) {
			tag := ""
			return makeRelease(tag, nil), nil
		}

		out, err := Run()

		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if out != "" {
			t.Errorf("got %q, want empty", out)
		}
	})

	t.Run("calls downloadAndReplace when newer version available", func(t *testing.T) {
		var downloadedURL string

		getLatestReleaseFn = func() (*github.RepositoryRelease, error) {
			assets := []*github.ReleaseAsset{
				makeAsset("sbx-darwin-arm64.tar.gz", "https://example.com/sbx-darwin-arm64.tar.gz"),
				makeAsset("sbx-darwin-amd64.tar.gz", "https://example.com/sbx-darwin-amd64.tar.gz"),
			}
			return makeRelease("v99.0.0", assets), nil
		}

		downloadAndReplaceFn = func(url string) error {
			downloadedURL = url
			return nil
		}

		out, err := Run()

		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if out != "updated!\n" {
			t.Errorf("got %q, want %q", out, "updated!\n")
		}
		if downloadedURL == "" {
			t.Error("downloadAndReplace was not called")
		}
	})

	t.Run("returns error when download fails", func(t *testing.T) {
		getLatestReleaseFn = func() (*github.RepositoryRelease, error) {
			assets := []*github.ReleaseAsset{
				makeAsset("sbx-darwin-arm64.tar.gz", "https://example.com/sbx-darwin-arm64.tar.gz"),
				makeAsset("sbx-darwin-amd64.tar.gz", "https://example.com/sbx-darwin-amd64.tar.gz"),
			}
			return makeRelease("v99.0.0", assets), nil
		}

		downloadAndReplaceFn = func(url string) error {
			return errors.New("permission denied")
		}

		_, err := Run()

		if err == nil {
			t.Error("got nil, want error")
		}
	})
}
