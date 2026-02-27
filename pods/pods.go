package pods

import (
	"fmt"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
)

const template = "https://app.datadoghq.com/orchestration/explorer/pod?query=kube_namespace:%s"

var openURL = open.Open

func Run() (string, error) {
	err := openURL(Url())

	if err != nil {
		return "", err
	}

	return "", nil
}

var nameFn = name.Name

func Url() string {
	name, _ := nameFn()
	return fmt.Sprintf(template, name)
}

var cmdFn = cli.Cmd

func GetPods(namespace, deployment string) ([]string, error) {
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
