package main

import (
	"errors"
	"log"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		_ = os.Remove("debug.log")

		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalln("fatal:", err)
		}
		defer f.Close()
	}

	configDir, err := getConfigDir()
	if err != nil {
		log.Fatalf("could not get config directory: %v", err)
	}
	cfgPath := configFilePath(configDir)
	fallbackEmail := getGlobalGitEmail()
	cfg, err := ensureConfig(cfgPath, fallbackEmail)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	dataDir, err := getDataPath()
	if err != nil {
		log.Fatalf("could not get data path: %v", err)
	}
	todayEntry, err := getTodayEntry(dataDir)
	if err != nil {
		log.Printf("could not load today entry: %v", err)
	}

	m := initialModel(dataDir, cfg, cfgPath, todayEntry)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}

func getTodayEntry(dataDir string) (*entry, error) {
	today := time.Now().Format(YYYY_MM_DD)
	filePath := entryFilePath(dataDir, today)
	e, err := loadEntry(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}
