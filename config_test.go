package main

import (
	"os"
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
