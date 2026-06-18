package main

import (
	"os/exec"
	"strings"
)

// getGlobalGitEmail returns the user's global Git email.
// This is the lowest-precedence fallback.
func getGlobalGitEmail() string {
	out, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// getRepoGitEmail returns the repo-local Git email only.
// It intentionally ignores the user's global Git config.
func getRepoGitEmail(repoPath string) string {
	out, err := exec.Command("git", "-C", repoPath, "config", "--local", "user.email").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// emailForRepo chooses the author email used when querying commits.
// Precedence: config repo email, repo-local Git email, global fallback.
func emailForRepo(repo repoConfig, fallbackEmail string) string {
	if repo.Email != "" {
		return repo.Email
	}

	if email := getRepoGitEmail(repo.Path); email != "" {
		return email
	}

	return fallbackEmail
}

func isInsideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
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
	email := emailForRepo(repo, fallbackEmail)
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
