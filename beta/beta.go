package beta

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/up"
)

const noBranch = "did not match any file(s) known to git"

const note = `»»»
Note: Navigate to the %s directory to view deployment details

`

var cmdFn = cli.Cmd
var upFn = up.Run
var repos = []string{"reverb", "search-remixer", "reverb-search-v2", "search-indexer"}

func Run() (string, error) {
	// Get current dir
	home, err := currentDir()

	if err != nil {
		return "", err
	}

	// Ensure home dir is in repos
	found := false
	for _, repo := range repos {
		if repo == home {
			found = true
			break
		}
	}

	if !found {
		repos = append(repos, home)
	}

	// get branch
	branch, err := getBranch()

	if err != nil {
		return "", err
	}

	// Deploy repos
	for _, repo := range repos {
		// change dir
		fmt.Printf("%s: ", repo)
		err = changeDir(repo)

		if err != nil {
			return "", fmt.Errorf("failed to change dir to %s: %w", repo, err)
		}

		// fetch
		err = fetch()

		if err != nil {
			return "", fmt.Errorf("failed to fetch repo %s: %w", repo, err)
		}

		// checkout branch
		_, err = checkoutBranch(branch)

		if err != nil {
			return "", fmt.Errorf("failed to checkout branch %s in repo %s: %w", branch, repo, err)
		}

		// sbx up
		out, err := upFn()

		if repo != home {
			formattedNote := fmt.Sprintf(note, repo)
			fmt.Println(formattedNote)
		}

		if err != nil {
			return out, fmt.Errorf("deploy failed for %s: %w", repo, err)
		}

		// for the beta, core migrations are needed to run the new search indexer predeploy task
		// five minutes give core time to build and run its predeploy migrations and seeding
		//
		// was getting multiple errors from flux patch step, three trying at once, so
		// added in a shorter sleep for the other repos to try to avoid this
		fmt.Println("waiting before continuing...")
		fmt.Println()
		if repo == "reverb" {
			sleep(240)
		} else {
			sleep(60)
		}
	}

	// go back to home dir
	fmt.Printf("returning to %s directory\n", home)
	err = changeDir(home)

	if err != nil {
		return "", fmt.Errorf("failed to change dir to %s: %w", home, err)
	}

	return "beta sandbox deploy triggered; monitor each deploy for completion", nil
}

func currentDir() (string, error) {
	gitRoot, err := cmdFn("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	out, err := cmdFn("basename", strings.TrimSpace(gitRoot))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func getBranch() (string, error) {
	out, err := cmdFn("git", "branch", "--show-current")

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func fetch() error {
	_, err := cmdFn("git", "fetch")

	if err != nil {
		return err
	}

	return nil
}

func checkoutBranch(branch string) (string, error) {
	out, err := cmdFn("git", "checkout", branch)

	if err != nil && strings.Contains(out, noBranch) {
		out, err = cmdFn("git", "checkout", "-b", branch, "main-sandbox")
		return out, err
	}

	return out, err
}

var changeDir = _changeDir

func _changeDir(dir string) error {
	return os.Chdir(fmt.Sprintf("../%s", dir))
}

var sleep = _sleep

func _sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
