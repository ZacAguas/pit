package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveEntryCmdReturnsSuccessMessage(t *testing.T) {
	dir := t.TempDir()

	e := entry{
		Date:     "2026-06-16",
		Did:      "did work",
		Blocked:  "blocked thing",
		Tomorrow: "next thing",
	}

	cmd := saveEntryCmd(dir, e)

	msg := cmd()
	got, ok := msg.(saveEntryMsg)
	if !ok {
		t.Fatalf("expected saveEntryMsg, got %T", msg)
	}
	if got.err != nil {
		t.Fatalf("expected no error, got %v", got.err)
	}

	loaded, err := loadEntry(entryFilePath(dir, e.Date))
	if err != nil {
		t.Fatal(err)
	}
	if loaded != e {
		t.Fatalf("expected %#v, got %#v", e, loaded)
	}
}

func TestTrackRepoCmdSavesConfigWithRepo(t *testing.T) {
	dir := t.TempDir()
	configPath := configFilePath(dir)
	repoPath := filepath.Join(dir, "project", ".")

	cmd := trackRepoCmd(configPath, config{}, repoPath)

	msg := cmd()
	got, ok := msg.(trackRepoMsg)
	if !ok {
		t.Fatalf("expected trackRepoMsg, got %T", msg)
	}
	if got.err != nil {
		t.Fatalf("expected no error, got %v", got.err)
	}

	loaded, err := loadConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}

	normalizedRepoPath, err := normalizeRepoPath(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	if !configHasRepo(loaded, normalizedRepoPath) {
		t.Fatalf("expected config to contain %q, got %#v", normalizedRepoPath, loaded.Repos)
	}
	if got.repoPath != normalizedRepoPath {
		t.Fatalf("expected message repo path %q, got %q", normalizedRepoPath, got.repoPath)
	}
}

func TestQueryReposCommitsCmdReturnsCommitsAndWarnings(t *testing.T) {
	dir := t.TempDir()
	repoDir := filepath.Join(dir, "repo")
	if err := os.Mkdir(repoDir, 0o755); err != nil {
		t.Fatal(err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", email)
	runGit(t, repoDir, "config", "user.name", name)

	filePath := filepath.Join(repoDir, file)
	if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	const commitMsg = "add test file"
	runGit(t, repoDir, "add", file)
	runGit(t, repoDir, "commit", "-m", commitMsg)

	cmd := queryReposCommitsCmd(
		[]repoConfig{
			{Path: repoDir},
			{Path: filepath.Join(dir, "missing")},
		},
		"2026-01-01",
		"fallback@example.com",
	)

	msg := cmd()
	got, ok := msg.(queryReposCommitsMsg)
	if !ok {
		t.Fatalf("expected queryReposCommitsMsg, got %T", msg)
	}

	wantCommits := "- " + commitMsg
	if got.commits != wantCommits {
		t.Fatalf("expected commits %q, got %q", wantCommits, got.commits)
	}
	if len(got.warnings) != 1 {
		t.Fatalf("expected one warning, got %#v", got.warnings)
	}
}
