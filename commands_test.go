package main

import (
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
