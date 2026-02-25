// ...existing code...
package console

import (
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("it successfully opens a console", func(t *testing.T) {
		nameFn = func() (string, error) {
			return "test-namespace", nil
		}
		checkVPNConnectionFn = func() error {
			return nil
		}
		checkClusterAccessFn = func() error {
			return nil
		}
		getPodsFn = func(namespace, deployment string) ([]string, error) {
			if namespace != "test-namespace" {
				t.Errorf("got namespace %s, want test-namespace", namespace)
			}
			if deployment != "reverb-web" {
				t.Errorf("got deployment %s, want reverb-web", deployment)
			}
			return []string{"pod-1"}, nil
		}
		originalExec := execConsole
		defer func() { execConsole = originalExec }()
		execConsole = func(namespace, pod, container string) (string, error) {
			if namespace != "test-namespace" {
				t.Errorf("got namespace %s, want test-namespace", namespace)
			}
			if pod != "pod-1" {
				t.Errorf("got pod %s, want pod-1", pod)
			}
			if container != "web-pumas" {
				t.Errorf("got container %s, want web-pumas", container)
			}
			return "console opened", nil
		}
		_, err := Run()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it errors when VPN is not connected", func(t *testing.T) {
		checkVPNConnectionFn = func() error {
			return errors.New("VPN connection check failed")
		}
		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "VPN connection check failed") {
			t.Errorf("expected error to contain 'VPN connection check failed', got %v", err.Error())
		}
	})

	t.Run("it errors when cluster access fails", func(t *testing.T) {
		checkVPNConnectionFn = func() error {
			return nil
		}
		checkClusterAccessFn = func() error {
			return errors.New("cluster access failed")
		}
		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "cluster access failed") {
			t.Errorf("expected error to contain 'cluster access failed', got %v", err.Error())
		}
	})

	t.Run("it errors when no pods found", func(t *testing.T) {
		checkVPNConnectionFn = func() error {
			return nil
		}
		checkClusterAccessFn = func() error {
			return nil
		}
		nameFn = func() (string, error) {
			return "test-namespace", nil
		}
		getPodsFn = func(namespace, deployment string) ([]string, error) {
			return []string{}, nil
		}
		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "no pods found") {
			t.Errorf("expected error to contain 'no pods found', got %v", err.Error())
		}
	})

	t.Run("it errors when getPods fails", func(t *testing.T) {
		checkVPNConnectionFn = func() error {
			return nil
		}
		checkClusterAccessFn = func() error {
			return nil
		}
		nameFn = func() (string, error) {
			return "test-namespace", nil
		}
		getPodsFn = func(namespace, deployment string) ([]string, error) {
			return nil, errors.New("getPods error")
		}
		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "failed to get pods") {
			t.Errorf("expected error to contain 'failed to get pods', got %v", err.Error())
		}
	})

	t.Run("it errors when execConsole fails", func(t *testing.T) {
		checkVPNConnectionFn = func() error {
			return nil
		}
		checkClusterAccessFn = func() error {
			return nil
		}
		nameFn = func() (string, error) {
			return "test-namespace", nil
		}
		getPodsFn = func(namespace, deployment string) ([]string, error) {
			return []string{"pod-1"}, nil
		}
		originalExec := execConsole
		defer func() { execConsole = originalExec }()
		execConsole = func(namespace, pod, container string) (string, error) {
			return "", errors.New("execConsole error")
		}
		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "execConsole error") {
			t.Errorf("expected error to contain 'execConsole error', got %v", err.Error())
		}
	})
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && (len(s) >= len(substr)) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
