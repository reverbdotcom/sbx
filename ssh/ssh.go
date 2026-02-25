package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/login"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/pods"
)

var cmdFn = cli.Cmd
var nameFn = name.Name
var execIntoContainer = _execIntoContainer
var selectItemFn = selectItem
var checkClusterAccessFn = login.CheckClusterAccess
var checkVPNConnectionFn = login.CheckVPNConnection
var loginFn = login.Run
var getPodsFn = pods.GetPods

const (
	defaultShell  = "/bin/sh"
	fallbackShell = "/bin/bash"
)

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

	// Get deployments
	deployments, err := getDeployments(namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get deployments: %w", err)
	}

	if len(deployments) == 0 {
		return "", fmt.Errorf("no deployments found in namespace %s", namespace)
	}

	// Select deployment
	deployment, err := selectItemFn("Select a deployment:", deployments)
	if err != nil {
		return "", fmt.Errorf("failed to select deployment: %w", err)
	}

	// Get pods for the selected deployment
	pods, err := getPodsFn(namespace, deployment)
	if err != nil {
		return "", fmt.Errorf("failed to get pods: %w", err)
	}

	if len(pods) == 0 {
		return "", fmt.Errorf("no pods found for deployment %s", deployment)
	}

	// Select pod (or use first if only one)
	var pod string
	if len(pods) == 1 {
		pod = pods[0]
	} else {
		pod, err = selectItemFn("Select a pod:", pods)
		if err != nil {
			return "", fmt.Errorf("failed to select pod: %w", err)
		}
	}

	// Get containers for the selected pod
	containers, err := getContainers(namespace, pod)
	if err != nil {
		return "", fmt.Errorf("failed to get containers: %w", err)
	}

	if len(containers) == 0 {
		return "", fmt.Errorf("no containers found in pod %s", pod)
	}

	// Select container (or use first if only one)
	var container string
	if len(containers) == 1 {
		container = containers[0]
	} else {
		container, err = selectItemFn("Select a container:", containers)
		if err != nil {
			return "", fmt.Errorf("failed to select container: %w", err)
		}
	}

	// Execute kubectl exec
	return execIntoContainer(namespace, pod, container)
}

func getDeployments(namespace string) ([]string, error) {
	out, err := cmdFn("kubectl", "get", "deployments", "-n", namespace, "-o", "jsonpath={.items[*].metadata.name}")
	if err != nil {
		return nil, fmt.Errorf("kubectl error: %s: %w", out, err)
	}

	deployments := strings.Fields(strings.TrimSpace(out))
	return deployments, nil
}

func getContainers(namespace, pod string) ([]string, error) {
	out, err := cmdFn("kubectl", "get", "pod", pod, "-n", namespace, "-o", "jsonpath={.spec.containers[*].name}")
	if err != nil {
		return nil, fmt.Errorf("kubectl error: %s: %w", out, err)
	}

	containers := strings.Fields(strings.TrimSpace(out))
	return containers, nil
}

func selectItem(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  50,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func buildShellCommand(shell string) string {
	// Validate shell is one of the expected values to prevent command injection
	if shell != defaultShell && shell != fallbackShell {
		// This should never happen since we only pass constants, but being defensive
		shell = defaultShell
	}
	// Command to start an interactive shell
	return "exec " + shell
}

func _execIntoContainer(namespace, pod, container string) (string, error) {
	// Using /bin/sh to run the conditional sourcing, then exec into the desired shell
	shellCmd := buildShellCommand(defaultShell)
	cmd := exec.Command("kubectl", "exec", "-it", "-n", namespace, pod, "-c", container, "--", "/bin/sh", "-c", shellCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// Try fallback shell if default shell fails
		shellCmd = buildShellCommand(fallbackShell)
		cmd = exec.Command("kubectl", "exec", "-it", "-n", namespace, pod, "-c", container, "--", "/bin/sh", "-c", shellCmd)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("failed to exec into container: %w", err)
		}
	}

	return "", nil
}
