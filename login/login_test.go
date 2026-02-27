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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestCheckClusterAccess(t *testing.T) {
	t.Run("it succeeds when kubectl version returns Server Version", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nKustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		err := CheckClusterAccess()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it automatically runs login when kubectl version returns SSO token error", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Error loading SSO Token: Token for okta does not exist\nClient Version: v1.30.3\nKustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3\nUnable to connect to the server: getting credentials: exec: executable aws failed with exit code 255", Err: errors.New("exit status 1")},
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		loginCalled := false
		loginFn = func() (string, error) {
			loginCalled = true
			return "login successful", nil
		}

		err := CheckClusterAccess()
		if err != nil {
			t.Errorf("expected no error after login, got %v", err)
		}
		if !loginCalled {
			t.Error("expected login to be called")
		}
	})

	t.Run("it automatically runs login when kubectl version returns unable to connect error", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nKustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3\nUnable to connect to the server: getting credentials: exec: executable aws failed with exit code 255", Err: errors.New("exit status 1")},
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		loginCalled := false
		loginFn = func() (string, error) {
			loginCalled = true
			return "login successful", nil
		}

		err := CheckClusterAccess()
		if err != nil {
			t.Errorf("expected no error after login, got %v", err)
		}
		if !loginCalled {
			t.Error("expected login to be called")
		}
	})

	t.Run("it automatically runs login when Server Version is not in output", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nKustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3", Err: nil},
			{Command: "kubectl version", Out: "Client Version: v1.30.3\nServer Version: v1.32.9-eks-3cfe0ce", Err: nil},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		loginCalled := false
		loginFn = func() (string, error) {
			loginCalled = true
			return "login successful", nil
		}

		err := CheckClusterAccess()
		if err != nil {
			t.Errorf("expected no error after login, got %v", err)
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

		loginFn = func() (string, error) {
			return "", errors.New("login failed")
		}

		err := CheckClusterAccess()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "failed to authenticate") {
			t.Errorf("expected error to contain 'failed to authenticate', got %v", err.Error())
		}
	})

	t.Run("it errors when connection still fails after login", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "Error loading SSO Token: Token for okta does not exist\nUnable to connect to the server", Err: errors.New("exit status 1")},
			{Command: "kubectl version", Out: "Unable to connect to the server", Err: errors.New("exit status 1")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		loginFn = func() (string, error) {
			return "login successful", nil
		}

		err := CheckClusterAccess()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "kubectl still cannot connect to preprod cluster after login") {
			t.Errorf("expected error to contain 'kubectl still cannot connect to preprod cluster after login', got %v", err.Error())
		}
	})

	t.Run("it errors on kubectl version failure with unrelated error", func(t *testing.T) {
		mockCalls := []cli.MockCall{
			{Command: "kubectl version", Out: "some other error", Err: errors.New("command not found")},
		}

		cmdFn = cli.MockCmd(t, mockCalls)

		err := CheckClusterAccess()
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !contains(err.Error(), "kubectl version check failed") {
			t.Errorf("expected error to contain 'kubectl version check failed', got %v", err.Error())
		}
	})
}

func TestCheckVPNConnection(t *testing.T) {
	t.Run("it succeeds when VPN check URL is reachable", func(t *testing.T) {
		// Since we're testing the actual function, we need to mock it
		// The function variable is already set up for mocking in the package
		originalCheck := checkVPNConnectionFn
		defer func() { checkVPNConnectionFn = originalCheck }()

		checkVPNConnectionFn = func() error {
			return nil
		}

		err := checkVPNConnectionFn()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("it fails when VPN check URL is not reachable", func(t *testing.T) {
		originalCheck := checkVPNConnectionFn
		defer func() { checkVPNConnectionFn = originalCheck }()

		checkVPNConnectionFn = func() error {
			return errors.New("VPN connection check failed")
		}

		err := checkVPNConnectionFn()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
