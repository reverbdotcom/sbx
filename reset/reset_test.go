package reset

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

const profileContent = "[profile preprod]\nsso_account_id = 123456789012\nsso_role_name = reverb-dev\n"

// restore captures the package-level seams and returns a func that resets them.
func restore() func() {
	origCmd := cmdFn
	origHome := homeDirFn
	origConfirm := confirmFn
	origFetch := fetchProfileFn
	origLogin := awsSSOLoginFn
	origCheck := checkCommandFn

	return func() {
		cmdFn = origCmd
		homeDirFn = origHome
		confirmFn = origConfirm
		fetchProfileFn = origFetch
		awsSSOLoginFn = origLogin
		checkCommandFn = origCheck
	}
}

func TestRun(t *testing.T) {
	defer restore()()

	t.Run("it resets the config and configures preprod", func(t *testing.T) {
		home := t.TempDir()

		checkCommandFn = func(string) error { return nil }
		homeDirFn = func() (string, error) { return home, nil }
		confirmFn = func() (bool, error) { return true, nil }
		fetchProfileFn = func() (string, error) { return profileContent, nil }
		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws eks update-kubeconfig --region us-east-1 --profile preprod --name preprod-v6 --alias preprod", Out: "Updated context preprod\n", Err: nil},
		})

		result, err := Run()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(result, "Reset complete") {
			t.Errorf("expected completion message, got %v", result)
		}

		got, err := os.ReadFile(filepath.Join(home, ".aws", "config"))
		if err != nil {
			t.Fatalf("expected AWS config to be written, got %v", err)
		}

		if string(got) != profileContent {
			t.Errorf("got %q, want %q", string(got), profileContent)
		}
	})

	t.Run("it backs up an existing AWS config", func(t *testing.T) {
		home := t.TempDir()
		awsDir := filepath.Join(home, ".aws")
		if err := os.MkdirAll(awsDir, dirPerm); err != nil {
			t.Fatal(err)
		}
		existing := "old config"
		if err := os.WriteFile(filepath.Join(awsDir, "config"), []byte(existing), filePerm); err != nil {
			t.Fatal(err)
		}

		checkCommandFn = func(string) error { return nil }
		homeDirFn = func() (string, error) { return home, nil }
		confirmFn = func() (bool, error) { return true, nil }
		fetchProfileFn = func() (string, error) { return profileContent, nil }
		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws eks update-kubeconfig --region us-east-1 --profile preprod --name preprod-v6 --alias preprod", Out: "Updated context preprod\n", Err: nil},
		})

		if _, err := Run(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		backup, err := os.ReadFile(filepath.Join(awsDir, "config.backup"))
		if err != nil {
			t.Fatalf("expected backup to be written, got %v", err)
		}

		if string(backup) != existing {
			t.Errorf("backup got %q, want %q", string(backup), existing)
		}
	})

	t.Run("it logs in and retries when the SSO token has expired", func(t *testing.T) {
		home := t.TempDir()

		checkCommandFn = func(string) error { return nil }
		homeDirFn = func() (string, error) { return home, nil }
		confirmFn = func() (bool, error) { return true, nil }
		fetchProfileFn = func() (string, error) { return profileContent, nil }
		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws eks update-kubeconfig --region us-east-1 --profile preprod --name preprod-v6 --alias preprod", Out: "Error when retrieving token from sso: Token has expired and refresh failed", Err: errors.New("exit status 255")},
			{Command: "aws eks update-kubeconfig --region us-east-1 --profile preprod --name preprod-v6 --alias preprod", Out: "Updated context preprod\n", Err: nil},
		})

		loginCalled := false
		awsSSOLoginFn = func(profile string) error {
			loginCalled = true
			return nil
		}

		if _, err := Run(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !loginCalled {
			t.Error("expected aws sso login to be called")
		}
	})

	t.Run("it errors when login fails after an expired token", func(t *testing.T) {
		home := t.TempDir()

		checkCommandFn = func(string) error { return nil }
		homeDirFn = func() (string, error) { return home, nil }
		confirmFn = func() (bool, error) { return true, nil }
		fetchProfileFn = func() (string, error) { return profileContent, nil }
		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws eks update-kubeconfig --region us-east-1 --profile preprod --name preprod-v6 --alias preprod", Out: "Error loading SSO Token", Err: errors.New("exit status 255")},
		})

		awsSSOLoginFn = func(profile string) error { return errors.New("browser closed") }

		_, err := Run()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "aws sso login failed") {
			t.Errorf("expected 'aws sso login failed', got %v", err.Error())
		}
	})

	t.Run("it aborts before touching files when the fetch fails", func(t *testing.T) {
		home := t.TempDir()
		awsDir := filepath.Join(home, ".aws")
		if err := os.MkdirAll(awsDir, dirPerm); err != nil {
			t.Fatal(err)
		}
		existing := "untouched"
		configFile := filepath.Join(awsDir, "config")
		if err := os.WriteFile(configFile, []byte(existing), filePerm); err != nil {
			t.Fatal(err)
		}

		checkCommandFn = func(string) error { return nil }
		homeDirFn = func() (string, error) { return home, nil }
		confirmFn = func() (bool, error) { return true, nil }
		fetchProfileFn = func() (string, error) { return "", errors.New("could not fetch AWS config template: 404") }
		cmdFn = cli.MockCmd(t, []cli.MockCall{})

		_, err := Run()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		got, _ := os.ReadFile(configFile)
		if string(got) != existing {
			t.Errorf("config should be untouched, got %q", string(got))
		}

		if _, err := os.Stat(configFile + ".backup"); !os.IsNotExist(err) {
			t.Error("no backup should be created when fetch fails")
		}
	})

	t.Run("it cancels when the user declines confirmation", func(t *testing.T) {
		home := t.TempDir()

		checkCommandFn = func(string) error { return nil }
		homeDirFn = func() (string, error) { return home, nil }
		confirmFn = func() (bool, error) { return false, nil }
		fetchCalled := false
		fetchProfileFn = func() (string, error) {
			fetchCalled = true
			return profileContent, nil
		}
		cmdFn = cli.MockCmd(t, []cli.MockCall{})

		result, err := Run()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if fetchCalled {
			t.Error("fetch should not run after declining")
		}

		if !strings.Contains(result, "cancelled") {
			t.Errorf("expected cancellation message, got %v", result)
		}
	})

	t.Run("it errors when aws CLI is not found", func(t *testing.T) {
		checkCommandFn = func(name string) error {
			if name == "aws" {
				return errors.New("not found")
			}
			return nil
		}

		_, err := Run()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "aws CLI is required") {
			t.Errorf("expected 'aws CLI is required', got %v", err.Error())
		}
	})
}

func TestCheckCommand(t *testing.T) {
	t.Run("it finds an existing command", func(t *testing.T) {
		if err := checkCommand("ls"); err != nil {
			t.Errorf("expected no error for 'ls', got %v", err)
		}
	})

	t.Run("it errors for a non-existent command", func(t *testing.T) {
		if err := checkCommand("this-command-does-not-exist-12345"); err == nil {
			t.Error("expected error for non-existent command, got nil")
		}
	})
}

func TestNeedsLogin(t *testing.T) {
	cases := []struct {
		out  string
		want bool
	}{
		{"Error loading SSO Token: Token for okta does not exist", true},
		{"Error when retrieving token from sso: Token has expired and refresh failed", true},
		{"The SSO session associated with this profile has expired or is otherwise invalid. To refresh this SSO session run aws sso login", true},
		{"Profile preprod does not exist", true},
		{"Updated context preprod", false},
		{"Cluster status is CREATING", false},
	}

	for _, c := range cases {
		if got := needsLogin(c.out); got != c.want {
			t.Errorf("needsLogin(%q) = %v, want %v", c.out, got, c.want)
		}
	}
}
