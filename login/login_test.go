package login

import (
	"errors"
	"strings"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestRun(t *testing.T) {
	// Save original functions
	origCheckCommandFn := checkCommandFn
	origAwsSSOLoginFn := awsSSOLoginFn

	// Restore after tests
	defer func() {
		checkCommandFn = origCheckCommandFn
		awsSSOLoginFn = origAwsSSOLoginFn
	}()

	t.Run("it successfully logs in and switches context", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "123456789012\n", Err: nil},
			{Command: "aws sts get-caller-identity --no-cli-pager", Out: `{"Account": "123456789012"}`, Err: nil},
			{Command: "kubectx preprod", Out: "Switched to context \"preprod\".\n", Err: nil},
		})

		result, err := Run()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		expectedResult := "Successfully logged in and switched to kubernetes context: preprod"
		if result != expectedResult {
			t.Errorf("got %v, want %v", result, expectedResult)
		}
	})

	t.Run("it logs in when session is expired", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		// Mock awsSSOLogin to succeed
		awsSSOLoginFn = func(profile string) error {
			return nil
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "123456789012\n", Err: nil},
			{Command: "aws sts get-caller-identity --no-cli-pager", Out: "", Err: errors.New("session expired")},
			{Command: "kubectx preprod", Out: "Switched to context \"preprod\".\n", Err: nil},
		})

		result, err := Run()

		if err != nil {
			t.Errorf("expected no error after successful login, got %v", err)
		}

		expectedResult := "Successfully logged in and switched to kubernetes context: preprod"
		if result != expectedResult {
			t.Errorf("got %v, want %v", result, expectedResult)
		}
	})

	t.Run("it returns error when account ID is empty", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "", Err: nil},
		})

		_, err := Run()

		if err == nil {
			t.Errorf("expected error for empty account ID, got nil")
		}

		expectedErr := "no sso_account_id found in AWS profile preprod"
		if err.Error() != expectedErr {
			t.Errorf("got %v, want %v", err.Error(), expectedErr)
		}
	})

	t.Run("it returns error when getting account ID fails", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "", Err: errors.New("profile not found")},
		})

		_, err := Run()

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to get AWS account ID") {
			t.Errorf("expected error message to contain 'failed to get AWS account ID', got %v", err.Error())
		}
	})

	t.Run("it returns error when kubectx fails", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "123456789012\n", Err: nil},
			{Command: "aws sts get-caller-identity --no-cli-pager", Out: `{"Account": "123456789012"}`, Err: nil},
			{Command: "kubectx preprod", Out: "context not found", Err: errors.New("exit status 1")},
		})

		_, err := Run()

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "failed to switch to kubernetes context") {
			t.Errorf("expected error message to contain 'failed to switch to kubernetes context', got %v", err.Error())
		}
	})

	t.Run("it handles different account ID", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		// Mock awsSSOLogin to succeed
		awsSSOLoginFn = func(profile string) error {
			return nil
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "123456789012\n", Err: nil},
			{Command: "aws sts get-caller-identity --no-cli-pager", Out: `{"Account": "999999999999"}`, Err: nil},
			{Command: "kubectx preprod", Out: "Switched to context \"preprod\".\n", Err: nil},
		})

		result, err := Run()

		if err != nil {
			t.Errorf("expected no error after successful login, got %v", err)
		}

		expectedResult := "Successfully logged in and switched to kubernetes context: preprod"
		if result != expectedResult {
			t.Errorf("got %v, want %v", result, expectedResult)
		}
	})

	t.Run("it returns error when aws command not found", func(t *testing.T) {
		// Mock checkCommand to fail for aws
		checkCommandFn = func(name string) error {
			if name == "aws" {
				return errors.New("not found")
			}
			return nil
		}

		_, err := Run()

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "aws CLI is required") {
			t.Errorf("expected error message to contain 'aws CLI is required', got %v", err.Error())
		}
	})

	t.Run("it returns error when kubectx command not found", func(t *testing.T) {
		// Mock checkCommand to fail for kubectx
		checkCommandFn = func(name string) error {
			if name == "kubectx" {
				return errors.New("not found")
			}
			return nil
		}

		_, err := Run()

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "kubectx is required") {
			t.Errorf("expected error message to contain 'kubectx is required', got %v", err.Error())
		}
	})

	t.Run("it returns error when aws sso login fails", func(t *testing.T) {
		// Mock checkCommand to always succeed
		checkCommandFn = func(name string) error {
			return nil
		}

		// Mock awsSSOLogin to fail
		awsSSOLoginFn = func(profile string) error {
			return errors.New("login failed")
		}

		cmdFn = cli.MockCmd(t, []cli.MockCall{
			{Command: "aws configure get sso_account_id --profile preprod", Out: "123456789012\n", Err: nil},
			{Command: "aws sts get-caller-identity --no-cli-pager", Out: "", Err: errors.New("session expired")},
		})

		_, err := Run()

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "aws sso login failed") {
			t.Errorf("expected error message to contain 'aws sso login failed', got %v", err.Error())
		}
	})
}

func TestCheckCommand(t *testing.T) {
	t.Run("it finds existing command", func(t *testing.T) {
		// Test with a command that should exist on any Unix-like system
		err := checkCommand("ls")

		if err != nil {
			t.Errorf("expected no error for 'ls' command, got %v", err)
		}
	})

	t.Run("it returns error for non-existent command", func(t *testing.T) {
		err := checkCommand("this-command-does-not-exist-12345")

		if err == nil {
			t.Errorf("expected error for non-existent command, got nil")
		}
	})
}
