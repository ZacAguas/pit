package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func entryFilePath(dir string, date string) string {
	const ext = ".json"
	return filepath.Join(dir, date+ext)
}

// Return an XDG compliant path for storing entries
func getDataPath() string {
	// TODO: use XDG_DATA_HOME

	// HACK: using temp dir for now
	baseDir := os.TempDir()

	dir := filepath.Join(baseDir, "pit")
	log.Printf("Using data directory %q", dir)
	return dir
}

func saveEntry(dir string, e entry) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	if e.Date == "" {
		return errors.New("cannot save entry with no date")
	}

	// create dir if it doesn't exist
	if err = os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	filePath := entryFilePath(dir, e.Date)
	return os.WriteFile(filePath, data, 0o644)
}

func loadEntry(filePath string) (entry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return entry{}, err
	}

	var e entry
	err = json.Unmarshal(data, &e)

	return e, err
}
