package main

import (
	"errors"
	"os/exec"
	"strings"
)

var errNoRepoGitEmail = errors.New("no repo git email")
var errNoGitEmail = errors.New("no git email configured")

// getGlobalGitEmail returns the user's global Git email.
// This is the lowest-precedence fallback.
func getGlobalGitEmail() (string, error) {
	out, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// getRepoGitEmail returns the repo-local Git email only.
// It intentionally ignores the user's global Git config.
func getRepoGitEmail(repoPath string) (string, error) {
	out, err := exec.Command("git", "-C", repoPath, "config", "--local", "user.email").Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return "", errNoRepoGitEmail
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// emailForRepo chooses the author email used when querying commits.
// Precedence: config repo email, repo-local Git email, global fallback.
func emailForRepo(repo repoConfig, fallbackEmail string) (string, error) {
	if repo.Email != "" {
		return repo.Email, nil
	}

	email, err := getRepoGitEmail(repo.Path)
	if err == nil && email != "" {
		return email, nil
	}
	if err != nil {
		if errors.Is(err, errNoRepoGitEmail) {
			return fallbackEmail, nil
		}
		return "", err
	}

	return fallbackEmail, nil
}

func isInsideGitRepo() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func currentRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func queryRepoCommits(repo repoConfig, sinceDate string, fallbackEmail string) (string, error) {
	email, err := emailForRepo(repo, fallbackEmail)
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", errNoGitEmail
	}
	// git -C <path> log --since=<date> --author=<email> --pretty=format:%s
	cmd := exec.Command("git", "-C", repo.Path, "log", "--since="+sinceDate, "--author="+email, "--pretty=format:%s")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return parseCommits(string(out)), nil
}

// Parse git log output into bullet-lines for did field
func parseCommits(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	commits := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		commits = append(commits, "- "+line)
	}
	return strings.Join(commits, "\n")
}
