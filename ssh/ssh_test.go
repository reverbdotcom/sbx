package ssh

import (
	"errors"
	"strings"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestRun(t *testing.T) {
	t.Run("it successfully drops into a container", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "app-1", Err: nil},
			{Command: "kubectl get deployment app-1 -n sandbox-test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "map[app:app-1]", Err: nil},
			{Command: "kubectl get pods -n sandbox-test-namespace -l app=app-1 -o jsonpath={.items[*].metadata.name}", Out: "pod-1", Err: nil},
			{Command: "kubectl get pod pod-1 -n sandbox-test-namespace -o jsonpath={.spec.containers[*].name}", Out: "container-1", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		// Mock VPN check to pass
		checkVPNConnectionFn = func() error {
			return nil
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

	t.Run("it errors when namespace cannot be determined", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		checkVPNConnectionFn = func() error {
			return nil
		}
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
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		checkVPNConnectionFn = func() error {
			return nil
		}
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
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "", Err: errors.New("kubectl error")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		checkVPNConnectionFn = func() error {
			return nil
		}
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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestRunWithClusterAccessCheck(t *testing.T) {
	t.Run("it checks cluster access before proceeding", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "app-1", Err: nil},
			{Command: "kubectl get deployment app-1 -n sandbox-test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "map[app:app-1]", Err: nil},
			{Command: "kubectl get pods -n sandbox-test-namespace -l app=app-1 -o jsonpath={.items[*].metadata.name}", Out: "pod-1", Err: nil},
			{Command: "kubectl get pod pod-1 -n sandbox-test-namespace -o jsonpath={.spec.containers[*].name}", Out: "container-1", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		checkVPNConnectionFn = func() error {
			return nil
		}
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		selectItemFn = func(label string, items []string) (string, error) {
			if len(items) > 0 {
				return items[0], nil
			}
			return "", errors.New("no items to select")
		}

		originalExec := execIntoContainer
		defer func() { execIntoContainer = originalExec }()

		execIntoContainer = func(namespace, pod, container string) (string, error) {
			return "", nil
		}

		_, err := Run()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it automatically logs in when cluster access check fails", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Error loading SSO Token: Token for okta does not exist\nUnable to connect to the server", Err: errors.New("exit status 1")},
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
			{Command: "kubectl get deployments -n sandbox-test-namespace -o jsonpath={.items[*].metadata.name}", Out: "app-1", Err: nil},
			{Command: "kubectl get deployment app-1 -n sandbox-test-namespace -o jsonpath={.spec.selector.matchLabels}", Out: "map[app:app-1]", Err: nil},
			{Command: "kubectl get pods -n sandbox-test-namespace -l app=app-1 -o jsonpath={.items[*].metadata.name}", Out: "pod-1", Err: nil},
			{Command: "kubectl get pod pod-1 -n sandbox-test-namespace -o jsonpath={.spec.containers[*].name}", Out: "container-1", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		checkVPNConnectionFn = func() error {
			return nil
		}
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		loginCalled := false
		loginFn = func() (string, error) {
			loginCalled = true
			return "login successful", nil
		}

		selectItemFn = func(label string, items []string) (string, error) {
			if len(items) > 0 {
				return items[0], nil
			}
			return "", errors.New("no items to select")
		}

		originalExec := execIntoContainer
		defer func() { execIntoContainer = originalExec }()

		execIntoContainer = func(namespace, pod, container string) (string, error) {
			return "", nil
		}

		_, err := Run()
		if err != nil {
			t.Errorf("expected no error after auto-login, got %v", err)
		}
		if !loginCalled {
			t.Error("expected login to be called")
		}
	})

	t.Run("it errors when login fails", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Error loading SSO Token: Token for okta does not exist\nUnable to connect to the server", Err: errors.New("exit status 1")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)
		checkVPNConnectionFn = func() error {
			return nil
		}
		nameFn = func() (string, error) {
			return "sandbox-test-namespace", nil
		}

		loginFn = func() (string, error) {
			return "", errors.New("login failed")
		}

		_, err := Run()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "failed to authenticate") {
			t.Errorf("expected error to contain 'failed to authenticate', got %v", err.Error())
		}
	})
}
