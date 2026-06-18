package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const (
	email = "repo@example.com"
	name  = "Test User"
	file  = "test"
)

func TestParseCommits(t *testing.T) {
	logOutput := "abc123 Add config\nabc456 Fix save\n"
	got := parseCommits(logOutput)
	want := "- abc123 Add config\n- abc456 Fix save"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParseCommitsIgnoresBlankLines(t *testing.T) {
	got := parseCommits("\nAdd config\n\nFix save\n")
	want := "- Add config\n- Fix save"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestParseCommitsEmptyOutput(t *testing.T) {
	got := parseCommits("")
	want := ""

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestQueryRepoCommits(t *testing.T) {
	dir := t.TempDir()

	// set-up repo
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", email)
	runGit(t, dir, "config", "user.name", name)

	// create test file
	_, err := os.Create(filepath.Join(dir, file))
	if err != nil {
		t.Fatal(err)
	}

	// commit test file
	const commitMsg = "initial commit"
	runGit(t, dir, "add", file)
	runGit(t, dir, "commit", "-m", commitMsg)

	got, err := queryRepoCommits(repoConfig{Path: dir}, "2026-01-01", "fallback@example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := "- " + commitMsg

	if want != got {
		t.Fatalf("expected: %q, got %q", want, got)
	}
}

func TestEmailForRepoUsesConfigEmailFirst(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "local@example.com")

	repo := repoConfig{
		Path:  dir,
		Email: "configured@example.com",
	}

	got, err := emailForRepo(repo, "fallback@example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := "configured@example.com"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestEmailForRepoUsesRepoLocalEmail(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "local@example.com")

	repo := repoConfig{Path: dir}

	got, err := emailForRepo(repo, "fallback@example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := "local@example.com"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestEmailForRepoFallsBackToGlobalFallback(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")

	repo := repoConfig{Path: dir}

	got, err := emailForRepo(repo, "fallback@example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := "fallback@example.com"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestGetRepoGitEmailReturnsNoRepoGitEmailError(t *testing.T) {
	dir := t.TempDir()
	runGit(t, dir, "init")

	_, err := getRepoGitEmail(dir)
	if !errors.Is(err, errNoRepoGitEmail) {
		t.Fatalf("expected errNoRepoGitEmail, got %v", err)
	}
}

func TestEmailForRepoReturnsUnexpectedGitError(t *testing.T) {
	_, err := emailForRepo(repoConfig{Path: filepath.Join(t.TempDir(), "missing")}, "fallback@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}
