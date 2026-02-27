package console

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/reverbdotcom/sbx/login"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/pods"
)

var nameFn = name.Name
var execConsole = _execConsole
var checkClusterAccessFn = login.CheckClusterAccess
var checkVPNConnectionFn = login.CheckVPNConnection
var getPodsFn = pods.GetPods

const (
	deployment = "reverb-web"
	container  = "web-pumas"
)

// Run opens a core rails console
func Run() (string, error) {
	// Check VPN connection
	if err := checkVPNConnectionFn(); err != nil {
		return "", err
	}

	// Check if kubectl can connect to the cluster
	if err := checkClusterAccessFn(); err != nil {
		return "", err
	}

	namespace, err := nameFn()
	if err != nil {
		return "", fmt.Errorf("failed to get namespace: %w", err)
	}

	// Get pods with label app=web-puma
	pods, err := getPodsFn(namespace, deployment)
	if err != nil {
		return "", fmt.Errorf("failed to get pods: %w", err)
	}

	if len(pods) == 0 {
		return "", fmt.Errorf("no pods found for deployment %s", deployment)
	}

	// Select first pod
	pod := pods[0]

	fmt.Printf("Running console on pod '%s' in sandbox '%s'\n\n", pod, namespace)

	return execConsole(namespace, pod, container)
}

func _execConsole(namespace, pod, container string) (string, error) {
	// Using /bin/sh to run the conditional sourcing, then exec into rails console
	railsCmd := "cd /app && exec bin/rails c"
	cmd := exec.Command("kubectl", "exec", "-it", "-n", namespace, pod, "-c", container, "--", "/bin/sh", "-c", railsCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to exec console: %w", err)
	}

	return "", nil
}
