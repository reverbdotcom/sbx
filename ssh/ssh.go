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
		return nil, err
	}

	deployments := strings.Fields(strings.TrimSpace(out))
	return deployments, nil
}

func getPods(namespace, deployment string) ([]string, error) {
	// First get the deployment's selector
	selectorOut, err := cmdFn("kubectl", "get", "deployment", deployment, "-n", namespace, "-o", "jsonpath={.spec.selector.matchLabels}")
	if err != nil {
		return nil, err
	}

	// If the selector is empty or we can't parse it, fallback to a simple label selector
	if strings.TrimSpace(selectorOut) == "" {
		// Try with common label patterns
		out, err := cmdFn("kubectl", "get", "pods", "-n", namespace, "-l", fmt.Sprintf("app=%s", deployment), "-o", "jsonpath={.items[*].metadata.name}")
		if err != nil {
			return nil, err
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
		return nil, err
	}

	pods := strings.Fields(strings.TrimSpace(out))
	return pods, nil
}

func parseSelector(selectorJSON string) string {
	// Simple parsing of kubectl jsonpath output for matchLabels
	// Input format: map[app:myapp version:v1]
	// Output format: app=myapp,version=v1

	// Remove "map[" prefix and "]" suffix
	selectorJSON = strings.TrimSpace(selectorJSON)
	selectorJSON = strings.TrimPrefix(selectorJSON, "map[")
	selectorJSON = strings.TrimSuffix(selectorJSON, "]")

	// Split by space and convert key:value to key=value
	pairs := strings.Fields(selectorJSON)
	for i, pair := range pairs {
		pairs[i] = strings.Replace(pair, ":", "=", 1)
	}

	return strings.Join(pairs, ",")
}

func getContainers(namespace, pod string) ([]string, error) {
	out, err := cmdFn("kubectl", "get", "pod", pod, "-n", namespace, "-o", "jsonpath={.spec.containers[*].name}")
	if err != nil {
		return nil, err
	}

	containers := strings.Fields(strings.TrimSpace(out))
	return containers, nil
}

func selectItem(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func _execIntoContainer(namespace, pod, container string) (string, error) {
	cmd := exec.Command("kubectl", "exec", "-it", "-n", namespace, pod, "-c", container, "--", "/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// Try /bin/bash if /bin/sh fails
		cmd = exec.Command("kubectl", "exec", "-it", "-n", namespace, pod, "-c", container, "--", "/bin/bash")
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
