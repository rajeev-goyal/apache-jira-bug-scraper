// In file: pkg/git/search.go
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// FindCommit searches a local git repo for a commit hash matching a bug ID.
// It runs: git log --all --grep="<bugID>" -n 1 --pretty=format:%H
//
// repoPath: The *absolute* or *relative* path to the .git directory (e.g., "../kafka")
// bugID:    The bug key to search for (e.g., "KAFKA-19112")
//
// Returns the commit hash as a string, or an empty string if not found.
func FindCommit(repoPath string, bugID string) (string, error) {
	// 1. Prepare the arguments for the 'git' command
	grepArg := fmt.Sprintf("--grep=%s", bugID)

	args := []string{
		"log",
		"--all",
		grepArg,
		"-n", "1", // Only find the first match
		"--pretty=format:%H", // Print *only* the commit hash
	}

	// 2. Set up the command
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath // This is crucial: run the command *inside* the repo path

	// 3. Run the command and capture STDOUT and STDERR
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()

	// 4. Handle errors
	if err != nil {
		// If 'git log' exits with an error, it might just mean "not found"
		// or it could be a real error. We'll check stderr.
		stderr := errbuf.String()
		if stderr != "" {
			return "", fmt.Errorf("git command failed: %s", stderr)
		}
		// If no error output, but exit code was non-zero, it likely means
		// 'git log' found no commits, which is not an *error* for us.
	}

	// 5. Clean and return the output (the commit hash)
	commitHash := strings.TrimSpace(outbuf.String())

	return commitHash, nil
}

// In file: pkg/git/search.go
// (Add this function *after* your existing FindCommit function)

// GetCommitDiff retrieves the diff for a specific commit, but only for
// files that match the 'test' patterns.
// It runs: git show <commitHash> -- "*Test.java" "*Test.scala"
func GetCommitDiff(repoPath string, commitHash string) (string, error) {
	// 1. Prepare the arguments for the 'git' command
	args := []string{
		"show",
		commitHash,
	}

	// 2. Set up the command
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	// 3. Run the command and capture STDOUT and STDERR
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()

	// 4. Handle errors
	if err != nil {
		stderr := errbuf.String()
		if stderr != "" {
			return "", fmt.Errorf("git command failed: %s", stderr)
		}
	}

	// 5. Return the diff text
	return outbuf.String(), nil
}
