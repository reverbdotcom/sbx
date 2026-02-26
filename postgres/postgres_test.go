package postgres

import (
	"errors"
	"testing"
)

func TestOpen(t *testing.T) {
	nameFn = func() (string, error) {
		return "test-sandbox", nil
	}

	openURL = func(url string) error {
		expectedUrl := "https://postgres-test-sandbox.int.orchestra.rvb.ai/"
		if url != expectedUrl {
			t.Errorf("Expected url %s, got %s", expectedUrl, url)
		}
		return nil
	}

	_, err := Open()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestOpen_Error(t *testing.T) {
	nameFn = func() (string, error) {
		return "test-sandbox", nil
	}

	openURL = func(url string) error {
		return errors.New("failed to open")
	}

	_, err := Open()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestOpenProgress(t *testing.T) {
	htmlUrlFn = func() (string, error) {
		return "https://github.com/test", nil
	}

	openURL = func(url string) error {
		expectedUrl := "https://github.com/test"
		if url != expectedUrl {
			t.Errorf("Expected url %s, got %s", expectedUrl, url)
		}
		return nil
	}

	_, err := OpenProgress()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
