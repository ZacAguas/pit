package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func entryFilePath(dir string, date string) string {
	const ext = ".json"
	return filepath.Join(dir, date+ext)
}

// Return an XDG compliant path for storing entries
func getDataPath() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(baseDir, "pit", "entries")
	return dir, nil
}

// Called in saveEntryCmd as it is a side-effect
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

func getAllEntries(dataDir string) ([]entry, error) {
	dirEntries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	var entries []entry
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".json") {
			continue
		}

		e, err := loadEntry(filepath.Join(dataDir, dirEntry.Name()))
		if err != nil {
			log.Printf("failed to load %q: %v", dirEntry.Name(), err)
			continue
		}

		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date > entries[j].Date
	})

	return entries, nil
}
