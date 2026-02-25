package login

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/reverbdotcom/sbx/cli"
)

const (
	profile     = "preprod"
	vpnCheckURL = "https://nsqadmin.reverb.tools/"
)

var cmdFn = cli.Cmd
var checkCommandFn = checkCommand
var awsSSOLoginFn = awsSSOLogin
var loginFn = Run
var checkVPNConnectionFn = CheckVPNConnection

// Run executes the k8s login workflow:
// 1. Check if aws and kubectx are available
// 2. Set AWS_PROFILE environment variable
// 3. Get account ID from AWS profile
// 4. Check if already authenticated
// 5. If not authenticated, run aws sso login
// 6. Switch to preprod kubernetes context using kubectx
func Run() (string, error) {
	// Check if required commands are available
	if err := checkCommandFn("aws"); err != nil {
		return "", fmt.Errorf("aws CLI is required but not found. Please install it: https://aws.amazon.com/cli/")
	}
	if err := checkCommandFn("kubectx"); err != nil {
		return "", fmt.Errorf("kubectx is required but not found. Please install it: https://github.com/ahmetb/kubectx")
	}

	// Set AWS_PROFILE environment variable
	os.Setenv("AWS_PROFILE", profile)

	// Get account ID from AWS profile
	accountID, err := cmdFn("aws", "configure", "get", "sso_account_id", "--profile", profile)
	if err != nil {
		return "", fmt.Errorf("failed to get AWS account ID from profile %s: %w", profile, err)
	}
	accountID = strings.TrimSpace(accountID)

	if accountID == "" {
		return "", fmt.Errorf("no sso_account_id found in AWS profile %s", profile)
	}

	// Check if already authenticated by calling aws sts get-caller-identity
	activeAccountID, err := cmdFn("aws", "sts", "get-caller-identity", "--no-cli-pager")
	if err != nil || !strings.Contains(activeAccountID, accountID) {
		// Not authenticated or wrong account, need to login
		fmt.Println("AWS session not active or expired. Logging in...")
		if err := awsSSOLoginFn(profile); err != nil {
			return "", fmt.Errorf("aws sso login failed: %w", err)
		}
		fmt.Println("AWS login successful.")
	} else {
		fmt.Println("AWS session is active.")
	}

	// Switch to preprod kubernetes context
	fmt.Printf("Switching to kubernetes context: %s...\n", profile)
	output, err := cmdFn("kubectx", profile)
	if err != nil {
		return "", fmt.Errorf("failed to switch to kubernetes context %s: %s: %w", profile, output, err)
	}

	return fmt.Sprintf("Successfully logged in and switched to kubernetes context: %s", profile), nil
}

// checkCommand checks if a command is available in PATH
func checkCommand(name string) error {
	_, err := exec.LookPath(name)
	return err
}

// awsSSOLogin runs the aws sso login command interactively
func awsSSOLogin(profile string) error {
	cmd := exec.Command("aws", "sso", "login", "--profile", profile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// checkClusterAccess checks if kubectl can connect to the preprod cluster
// If not authenticated, it automatically runs sbx k8s login
func CheckClusterAccess() error {
	out, err := cmdFn("kubectl", "version")
	if err != nil {
		// Check if the error is due to SSO token or connection issues
		if strings.Contains(out, "Error loading SSO Token") ||
			strings.Contains(out, "Unable to connect to the server") ||
			strings.Contains(out, "getting credentials") {
			return attemptAutoLogin()
		}
		return fmt.Errorf("kubectl version check failed: %s: %w", out, err)
	}

	// Verify that Server Version is present in the output
	if !strings.Contains(out, "Server Version") {
		return attemptAutoLogin()
	}

	return nil
}

// attemptAutoLogin attempts to authenticate and verify connection
func attemptAutoLogin() error {
	fmt.Println("kubectl cannot connect to preprod cluster. Running 'sbx k8s login'...")

	// Attempt to login
	_, loginErr := loginFn()
	if loginErr != nil {
		return fmt.Errorf("failed to authenticate: %w", loginErr)
	}

	// Verify connection after login
	out, err := cmdFn("kubectl", "version")
	if err != nil || !strings.Contains(out, "Server Version") {
		return fmt.Errorf("kubectl still cannot connect to preprod cluster after login")
	}

	return nil
}

// CheckVPNConnection verifies that the VPN is connected by making a request to nsqadmin
func CheckVPNConnection() error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(vpnCheckURL)
	if err != nil {
		return fmt.Errorf("VPN connection check failed. Please ensure you are connected to the VPN.\nError: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("VPN connection check failed with status %d. Please ensure you are connected to the VPN", resp.StatusCode)
	}

	return nil
}
