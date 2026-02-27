package pods

import (
	"errors"
	"fmt"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestUrl(t *testing.T) {
	name := "sandbox-foo-bar-baz"
	nameFn = func() (string, error) { return name, nil }

	t.Run("it returns url", func(t *testing.T) {
		got := Url()
		want := fmt.Sprintf(template, name)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestRun(t *testing.T) {
	name := "sandbox-foo-bar-baz"
	nameFn = func() (string, error) { return name, nil }

	t.Run("it opens url", func(t *testing.T) {
		want := fmt.Sprintf(template, name)

		openURL = func(got string) error {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}

			return nil
		}

		_, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it returns err", func(t *testing.T) {
		openURL = func(_ string) error {
			return errors.New("open error")
		}

		_, err := Run()

		want := "open error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
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

		pods, err := GetPods("test-namespace", "app-1")
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

		pods, err := GetPods("test-namespace", "app-1")
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

		_, err := GetPods("test-namespace", "app-1")
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
		{
			name:     "label with dots",
			input:    `map["reverb.com/deployment":"graphql-gateway"]`,
			expected: "reverb.com/deployment=graphql-gateway",
		},
		{
			name:     "multiple labels with special characters",
			input:    `map["reverb.com/deployment":"graphql-gateway" "app":"web"]`,
			expected: "reverb.com/deployment=graphql-gateway,app=web",
		},
		{
			name:     "mixed labels with and without quotes",
			input:    `map["reverb.com/deployment":"graphql-gateway" version:v1]`,
			expected: "reverb.com/deployment=graphql-gateway,version=v1",
		},
		{
			name:     "JSON format single label",
			input:    `{"reverb.com/deployment":"graphql-gateway"}`,
			expected: "reverb.com/deployment=graphql-gateway",
		},
		{
			name:     "JSON format multiple labels",
			input:    `{"reverb.com/deployment":"graphql-gateway","app":"web"}`,
			expected: "reverb.com/deployment=graphql-gateway,app=web",
		},
		{
			name:     "JSON format empty",
			input:    `{}`,
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
