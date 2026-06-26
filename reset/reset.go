package reset

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/github"
)

const (
	profile     = "preprod"
	cluster     = "preprod-v6"
	region      = "us-east-1"
	profilePath = "setup/aws-dev.profile"
	dirPerm     = 0755
	filePerm    = 0644
)

var cmdFn = cli.Cmd
var homeDirFn = os.UserHomeDir
var confirmFn = confirm
var fetchProfileFn = fetchProfile
var awsSSOLoginFn = awsSSOLogin
var checkCommandFn = checkCommand

// Run resets the user's AWS and kubernetes configuration so that
// `sbx k8s login` can authenticate against the preprod cluster.
//
// It fetches the AWS config template from the internal k8x repository,
// backs up any existing config, writes the fresh config, and rebuilds
// the preprod kubeconfig context.
func Run() (string, error) {
	if err := checkCommandFn("aws"); err != nil {
		return "", fmt.Errorf("aws CLI is required but not found. Please install it: https://aws.amazon.com/cli/")
	}

	homeDir, err := homeDirFn()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	fmt.Println("⚠️  This will override your AWS and kubernetes configurations!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ok, err := confirmFn()
	if err != nil {
		return "", err
	}
	if !ok {
		return "Reset cancelled.", nil
	}

	// Fetch the template before touching any files so a failed fetch
	// never destroys an existing configuration.
	content, err := fetchProfileFn()
	if err != nil {
		return "", err
	}

	if err := writeAWSConfig(homeDir, content); err != nil {
		return "", err
	}

	awsConfig := filepath.Join(homeDir, ".aws", "config")
	fmt.Printf("\n✓ AWS config reset\n")
	fmt.Printf("  Config written to: %s\n", awsConfig)

	if err := resetKubeconfig(homeDir); err != nil {
		return "", err
	}

	fmt.Println("\n⚙️  Setting up preprod kubeconfig...")
	if err := setupPreprodKubeconfig(); err != nil {
		return "", err
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	return "✓ Reset complete! Run 'sbx k8s login' to authenticate and switch to the preprod context.", nil
}

// fetchProfile retrieves the AWS config template from the internal k8x repository.
func fetchProfile() (string, error) {
	content, err := github.GetFileContents("k8x", profilePath)
	if err != nil {
		return "", fmt.Errorf("could not fetch AWS config template: %w", err)
	}

	return content, nil
}

// writeAWSConfig backs up any existing ~/.aws/config and writes the fetched content.
func writeAWSConfig(homeDir, content string) error {
	awsDir := filepath.Join(homeDir, ".aws")
	if err := os.MkdirAll(awsDir, dirPerm); err != nil {
		return fmt.Errorf("failed to create .aws directory: %w", err)
	}

	configFile := filepath.Join(awsDir, "config")
	if err := backup(configFile); err != nil {
		return err
	}

	if err := os.WriteFile(configFile, []byte(content), filePerm); err != nil {
		return fmt.Errorf("failed to write AWS config file: %w", err)
	}

	return nil
}

// resetKubeconfig backs up any existing ~/.kube/config so it can be rebuilt cleanly.
func resetKubeconfig(homeDir string) error {
	kubeDir := filepath.Join(homeDir, ".kube")
	if err := os.MkdirAll(kubeDir, dirPerm); err != nil {
		return fmt.Errorf("failed to create .kube directory: %w", err)
	}

	if err := backup(filepath.Join(kubeDir, "config")); err != nil {
		return err
	}

	return nil
}

// backup moves an existing file to "<path>.backup". It is a no-op if the file
// does not exist.
func backup(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	backupPath := path + ".backup"
	if err := os.Rename(path, backupPath); err != nil {
		return fmt.Errorf("failed to back up %s: %w", path, err)
	}

	fmt.Printf("  Backed up %s to %s\n", path, backupPath)
	return nil
}

// setupPreprodKubeconfig configures the preprod EKS context, triggering an SSO
// login and retrying once if the AWS session is missing or expired.
func setupPreprodKubeconfig() error {
	fmt.Printf("  Configuring %s cluster...\n", profile)

	out, err := updateKubeconfig()
	if err != nil && needsLogin(out) {
		fmt.Printf("    🔑 AWS session missing or expired. Logging in to %s...\n", profile)
		if loginErr := awsSSOLoginFn(profile); loginErr != nil {
			return fmt.Errorf("aws sso login failed for %s: %w", profile, loginErr)
		}

		out, err = updateKubeconfig()
	}

	if err != nil {
		return fmt.Errorf("failed to configure %s cluster: %s: %w", profile, strings.TrimSpace(out), err)
	}

	fmt.Printf("    ✓ %s configured\n", profile)
	return nil
}

func updateKubeconfig() (string, error) {
	return cmdFn("aws", "eks", "update-kubeconfig",
		"--region", region,
		"--profile", profile,
		"--name", cluster,
		"--alias", profile)
}

// needsLogin reports whether the aws output indicates the SSO session is
// missing, expired, or otherwise invalid and an `aws sso login` is required.
func needsLogin(out string) bool {
	lower := strings.ToLower(out)
	signals := []string{
		"sso token",         // Error loading SSO Token
		"token from sso",    // Error when retrieving token from sso
		"token has expired", // Token has expired and refresh failed
		"refresh failed",
		"sso session", // SSO session ... has expired or is otherwise invalid
		"sso login",   // ...run aws sso login with the corresponding profile
		"does not exist",
	}

	for _, s := range signals {
		if strings.Contains(lower, s) {
			return true
		}
	}

	return false
}

// confirm prompts the user to confirm the destructive reset.
func confirm() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("This overwrites your ~/.aws/config and ~/.kube/config. Continue? (y/N): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	answer := strings.TrimSpace(strings.ToLower(input))
	return answer == "y" || answer == "yes", nil
}

// awsSSOLogin runs the aws sso login command interactively.
func awsSSOLogin(profile string) error {
	cmd := exec.Command("aws", "sso", "login", "--profile", profile, "--region", region)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// checkCommand checks if a command is available in PATH.
func checkCommand(name string) error {
	_, err := exec.LookPath(name)
	return err
}
