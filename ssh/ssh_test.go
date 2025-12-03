package ssh

import (
	"errors"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestRun(t *testing.T) {
	t.Run("it successfully drops into a container", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "app-1", Err: nil},
			{Command: "kubectl get deployment app-1 -n sandbox-test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "map[app:app-1]", Err: nil},
			{Command: "kubectl get pods -n sandbox-test-namespace -l app=app-1 -o jsonpath={.items[*].metadata.name}", Out: "pod-1", Err: nil},
			{Command: "kubectl get pod pod-1 -n sandbox-test-namespace -o jsonpath={.spec.containers[*].name}", Out: "container-1", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		// Mock the select function to avoid interactive prompt
		selectItemFn = func(label string, items []string) (string, error) {
			if len(items) > 0 {
				return items[0], nil
			}
			return "", errors.New("no items to select")
		}

		// Mock the exec function to avoid actually running kubectl exec
		originalExec := execIntoContainer
		defer func() { execIntoContainer = originalExec }()

		execIntoContainer = func(namespace, pod, container string) (string, error) {
			if namespace != "sandbox-test-namespace" {
				t.Errorf("got namespace %s, want sandbox-test-namespace", namespace)
			}
			if pod != "pod-1" {
				t.Errorf("got pod %s, want pod-1", pod)
			}
			if container != "container-1" {
				t.Errorf("got container %s, want container-1", container)
			}
			return "", nil
		}

		_, err := Run()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it errors when namespace cannot be determined", func(t *testing.T) {
		nameFn = func() (string, error) {
			return "", errors.New("namespace error")
		}

		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "failed to get namespace") {
			t.Errorf("expected error to contain 'failed to get namespace', got %v", err.Error())
		}
	})

	t.Run("it errors when no deployments found", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "no deployments found") {
			t.Errorf("expected error to contain 'no deployments found', got %v", err.Error())
		}
	})

	t.Run("it errors when kubectl get deployments fails", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "", Err: errors.New("kubectl error")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "failed to get deployments") {
			t.Errorf("expected error to contain 'failed to get deployments', got %v", err.Error())
		}
	})
}

func TestGetDeployments(t *testing.T) {
	t.Run("it returns deployments", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployments -n test-namespace -o jsonpath={.items[*].metadata.name}", Out: "app-1 app-2 app-3", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		deployments, err := getDeployments("test-namespace")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(deployments) != 3 {
			t.Errorf("expected 3 deployments, got %d", len(deployments))
		}

		if deployments[0] != "app-1" || deployments[1] != "app-2" || deployments[2] != "app-3" {
			t.Errorf("unexpected deployment names: %v", deployments)
		}
	})

	t.Run("it returns empty list when no deployments", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployments -n test-namespace -o jsonpath={.items[*].metadata.name}", Out: "", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		deployments, err := getDeployments("test-namespace")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(deployments) != 0 {
			t.Errorf("expected 0 deployments, got %d", len(deployments))
		}
	})

	t.Run("it returns error on kubectl failure", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployments -n test-namespace -o jsonpath={.items[*].metadata.name}", Out: "", Err: errors.New("kubectl error")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		_, err := getDeployments("test-namespace")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetPods(t *testing.T) {
	t.Run("it returns pods for deployment using selector", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployment app-1 -n test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "map[app:app-1 version:v1]", Err: nil},
			{Command: "kubectl get pods -n test-namespace -l app=app-1,version=v1 -o jsonpath={.items[*].metadata.name}", Out: "pod-1 pod-2", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		pods, err := getPods("test-namespace", "app-1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(pods) != 2 {
			t.Errorf("expected 2 pods, got %d", len(pods))
		}

		if pods[0] != "pod-1" || pods[1] != "pod-2" {
			t.Errorf("unexpected pod names: %v", pods)
		}
	})

	t.Run("it falls back to app label when selector is empty", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployment app-1 -n test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "", Err: nil},
			{Command: "kubectl get pods -n test-namespace -l app=app-1 -o jsonpath={.items[*].metadata.name}", Out: "pod-1", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		pods, err := getPods("test-namespace", "app-1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(pods) != 1 {
			t.Errorf("expected 1 pod, got %d", len(pods))
		}

		if pods[0] != "pod-1" {
			t.Errorf("unexpected pod name: %v", pods)
		}
	})

	t.Run("it returns error on kubectl failure", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get deployment app-1 -n test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "", Err: errors.New("kubectl error")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		_, err := getPods("test-namespace", "app-1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetContainers(t *testing.T) {
	t.Run("it returns containers for pod", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get pod pod-1 -n test-namespace -o jsonpath={.spec.containers[*].name}", Out: "container-1 container-2", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		containers, err := getContainers("test-namespace", "pod-1")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if len(containers) != 2 {
			t.Errorf("expected 2 containers, got %d", len(containers))
		}

		if containers[0] != "container-1" || containers[1] != "container-2" {
			t.Errorf("unexpected container names: %v", containers)
		}
	})

	t.Run("it returns error on kubectl failure", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl get pod pod-1 -n test-namespace -o jsonpath={.spec.containers[*].name}", Out: "", Err: errors.New("kubectl error")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		_, err := getContainers("test-namespace", "pod-1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestParseSelector(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single label",
			input:    "map[app:myapp]",
			expected: "app=myapp",
		},
		{
			name:     "multiple labels",
			input:    "map[app:myapp version:v1]",
			expected: "app=myapp,version=v1",
		},
		{
			name:     "multiple labels with different order",
			input:    "map[version:v1 app:myapp tier:backend]",
			expected: "version=v1,app=myapp,tier=backend",
		},
		{
			name:     "empty map",
			input:    "map[]",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSelector(tt.input)
			if result != tt.expected {
				t.Errorf("parseSelector(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
