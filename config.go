package main

import (
	"errors"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
)

type config struct {
	GlobalEmail string       `toml:"global_email"`
	Repos       []repoConfig `toml:"repos"`
}

type repoConfig struct {
	Path  string `toml:"path"`
	Email string `toml:"email,omitempty"`
}

func getConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(dir, "pit")
	return configDir, nil
}

func configFilePath(configDir string) string {
	return filepath.Join(configDir, "config.toml")
}

func loadConfig(path string) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}
	var cfg config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return config{}, err
	}
	return cfg, nil
}

func saveConfig(filePath string, cfg config) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0o644)
}

// Load or create config
func ensureConfig(filePath string, fallbackEmail string) (config, error) {
	cfg, err := loadConfig(filePath)
	// successful load
	if err == nil {
		return cfg, nil
	}
	// unknown error
	if !errors.Is(err, os.ErrNotExist) {
		return config{}, err
	}

	// config doesn't exist
	cfg = config{
		GlobalEmail: fallbackEmail,
		Repos:       nil,
	}
	if err := saveConfig(filePath, cfg); err != nil {
		return config{}, err
	}
	return cfg, nil
}

func normalizeRepoPath(path string) (string, error) {
	return filepath.Abs(filepath.Clean(path))
}

func configHasRepo(cfg config, path string) bool {
	normalizedPath, err := normalizeRepoPath(path)
	if err != nil {
		return false
	}

	return slices.ContainsFunc(cfg.Repos, func(r repoConfig) bool {
		repoPath, err := normalizeRepoPath(r.Path)
		return err == nil && repoPath == normalizedPath
	})
}

func configWithRepo(cfg config, path string) (config, string, error) {
	repoPath, err := normalizeRepoPath(path)
	if err != nil {
		return config{}, "", err
	}

	if configHasRepo(cfg, repoPath) {
		return cfg, repoPath, nil
	}

	cfg.Repos = append(cfg.Repos, repoConfig{Path: repoPath})
	return cfg, repoPath, nil
}
