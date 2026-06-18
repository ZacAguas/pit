package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const fallback = "fallback@gmail.com"

func TestConfigCreatedWithFallbackEmail(t *testing.T) {
	dir := t.TempDir()
	filePath := configFilePath(dir)
	cfg, err := ensureConfig(filePath, fallback)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GlobalEmail != fallback {
		t.Fatalf("expected %q, got %q", fallback, cfg.GlobalEmail)
	}
}

func TestExistingConfigNotOverwritten(t *testing.T) {
	existing := config{
		GlobalEmail: "global@gmail.com",
		Repos: []repoConfig{
			{
				Path:  "a",
				Email: "local@gmail.com",
			},
		},
	}

	dir := t.TempDir()
	filePath := configFilePath(dir)
	if err := saveConfig(filePath, existing); err != nil {
		t.Fatal(err)
	}
	cfg, err := ensureConfig(filePath, fallback)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg, existing) {
		t.Fatalf("expected %#v, got %#v", existing, cfg)
	}
}

func TestMalformedTOMLReturnsError(t *testing.T) {
	malformed := `
GlobalEmail = "admin@example.com"

# ERROR: Mixed types inside the array
Repos = ["frontend-repo", 42, { Name = "backend" }]
	`

	dir := t.TempDir()
	filePath := configFilePath(dir)
	// write toml directly to dir
	err := os.WriteFile(filePath, []byte(malformed), 0o744)
	if err != nil {
		t.Fatalf("error writing malformed toml file: %v", err)
	}

	cfg, err := loadConfig(filePath)
	// should throw error
	if err == nil {
		t.Fatalf("expected error, got %#v", cfg)
	}
}

func TestSaveConfigWritesToPassedTempPath(t *testing.T) {
	dir := t.TempDir()
	filePath := configFilePath(dir)

	want := config{
		GlobalEmail: "global@gmail.com",
		Repos: []repoConfig{
			{
				Path:  "/tmp/project",
				Email: "repo@gmail.com",
			},
		},
	}

	if err := saveConfig(filePath, want); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected config file at %q: %v", filePath, err)
	}

	got, err := loadConfig(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

func TestConfigHasRepoFindsRepo(t *testing.T) {
	repoPath := filepath.Join(t.TempDir(), "project")

	cfg := config{
		Repos: []repoConfig{
			{
				Path: repoPath,
			},
		},
	}

	if !configHasRepo(cfg, repoPath) {
		t.Fatal("expected true, got false")
	}
}

func TestConfigHasRepoReturnsFalseForUnknownRepo(t *testing.T) {
	const unknownRepoPath = "/tmp/unknown"

	cfg := config{
		Repos: []repoConfig{
			{
				Path: "/tmp/project",
			},
		},
	}

	if configHasRepo(cfg, unknownRepoPath) {
		t.Fatal("expected false, got true")
	}
}

func TestConfigHasRepoNormalizesPaths(t *testing.T) {
	dir := t.TempDir()
	repoPath := filepath.Join(dir, "project")

	cfg := config{
		Repos: []repoConfig{
			{
				Path: filepath.Join(repoPath, "."),
			},
		},
	}

	if !configHasRepo(cfg, repoPath) {
		t.Fatal("expected true, got false")
	}
}

func TestConfigWithRepoAddsNormalizedPath(t *testing.T) {
	dir := t.TempDir()
	repoPath := filepath.Join(dir, "project", ".")

	got, normalizedRepoPath, err := configWithRepo(config{}, repoPath)
	if err != nil {
		t.Fatal(err)
	}

	wantPath, err := normalizeRepoPath(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	if normalizedRepoPath != wantPath {
		t.Fatalf("expected normalized path %q, got %q", wantPath, normalizedRepoPath)
	}
	if len(got.Repos) != 1 || got.Repos[0].Path != wantPath {
		t.Fatalf("expected repo path %q, got %#v", wantPath, got.Repos)
	}
}

func TestConfigWithRepoDoesNotDuplicateEquivalentPath(t *testing.T) {
	dir := t.TempDir()
	repoPath := filepath.Join(dir, "project")

	cfg := config{
		Repos: []repoConfig{{Path: filepath.Join(repoPath, ".")}},
	}

	got, _, err := configWithRepo(cfg, repoPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Repos) != 1 {
		t.Fatalf("expected one repo, got %#v", got.Repos)
	}
}
