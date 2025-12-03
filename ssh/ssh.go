package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
)

var cmdFn = cli.Cmd
var nameFn = name.Name
var execIntoContainer = _execIntoContainer
var selectItemFn = selectItem

const (
	defaultShell  = "/bin/sh"
	fallbackShell = "/bin/bash"
)

func Run() (string, error) {
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
	pods, err := getPods(namespace, deployment)
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

func getPods(namespace, deployment string) ([]string, error) {
	// First get the deployment's selector
	selectorOut, err := cmdFn("kubectl", "get", "deployment", deployment, "-n", namespace, "-o", "jsonpath={.spec.selector.matchLabels}")
	if err != nil {
		return nil, fmt.Errorf("kubectl error: %s: %w", selectorOut, err)
	}

	// If the selector is empty or we can't parse it, fallback to a simple label selector
	if strings.TrimSpace(selectorOut) == "" {
		// Try with common label patterns
		out, err := cmdFn("kubectl", "get", "pods", "-n", namespace, "-l", fmt.Sprintf("app=%s", deployment), "-o", "jsonpath={.items[*].metadata.name}")
		if err != nil {
			return nil, fmt.Errorf("kubectl error: %s: %w", out, err)
		}
		pods := strings.Fields(strings.TrimSpace(out))
		return pods, nil
	}

	// Get pods using the deployment's label selector
	// The selector output is in format: map[key1:value1 key2:value2]
	// We need to convert it to kubectl label selector format: key1=value1,key2=value2
	selector := parseSelector(selectorOut)

	out, err := cmdFn("kubectl", "get", "pods", "-n", namespace, "-l", selector, "-o", "jsonpath={.items[*].metadata.name}")
	if err != nil {
		return nil, fmt.Errorf("kubectl error: %s: %w", out, err)
	}

	pods := strings.Fields(strings.TrimSpace(out))
	return pods, nil
}

func parseSelector(selectorJSON string) string {
	// Simple parsing of kubectl jsonpath output for matchLabels
	// Input format can be:
	//   map[app:myapp version:v1]
	//   map["reverb.com/deployment":"graphql-gateway"]
	//   {"reverb.com/deployment":"graphql-gateway"}  (JSON format)
	// Output format: app=myapp,version=v1 or reverb.com/deployment=graphql-gateway

	selectorJSON = strings.TrimSpace(selectorJSON)

	// Handle JSON format (starts with {)
	if strings.HasPrefix(selectorJSON, "{") {
		// Remove { and }
		selectorJSON = strings.TrimPrefix(selectorJSON, "{")
		selectorJSON = strings.TrimSuffix(selectorJSON, "}")
		selectorJSON = strings.TrimSpace(selectorJSON)

		// Handle empty JSON object
		if selectorJSON == "" {
			return ""
		}

		// Parse JSON-style key:value pairs
		// Split by comma first (for multiple labels in JSON)
		parts := strings.Split(selectorJSON, ",")
		result := []string{}
		for _, part := range parts {
			part = strings.TrimSpace(part)
			// Remove quotes and convert : to =
			part = strings.ReplaceAll(part, "\"", "")
			part = strings.Replace(part, ":", "=", 1)
			if part != "" {
				result = append(result, part)
			}
		}
		return strings.Join(result, ",")
	}

	// Remove "map[" prefix and "]" suffix for Go map format
	selectorJSON = strings.TrimPrefix(selectorJSON, "map[")
	selectorJSON = strings.TrimSuffix(selectorJSON, "]")

	// Handle labels with special characters (quoted format)
	// Example: "reverb.com/deployment":"graphql-gateway" "app":"web"
	if strings.Contains(selectorJSON, "\"") {
		// Split by space to get individual key:value pairs
		pairs := strings.Fields(selectorJSON)
		result := []string{}
		for _, pair := range pairs {
			// Remove quotes and convert : to =
			pair = strings.ReplaceAll(pair, "\"", "")
			pair = strings.Replace(pair, ":", "=", 1)
			result = append(result, pair)
		}
		return strings.Join(result, ",")
	}

	// Handle simple labels without special characters
	// Example: app:myapp version:v1
	pairs := strings.Fields(selectorJSON)
	for i, pair := range pairs {
		pairs[i] = strings.Replace(pair, ":", "=", 1)
	}

	return strings.Join(pairs, ",")
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
	// Command to check if /etc/secrets/env exists and source it, then start an interactive shell
	return "[ -f /etc/secrets/env ] && . /etc/secrets/env; exec " + shell
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
