package main

import (
	"errors"
	"flag"
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

	daysBack := flag.Int("days-back", 1, "number of workdays back to load commits from")
	flag.Parse()

	configDir, err := getConfigDir()
	if err != nil {
		log.Fatalf("could not get config directory: %v", err)
	}
	cfgPath := configFilePath(configDir)
	fallbackEmail, err := getGlobalGitEmail()
	if err != nil {
		log.Printf("could not get global git email: %v", err)
	}
	cfg, err := ensureConfig(cfgPath, fallbackEmail)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// get tracked/untracked git repos
	var untrackedRepoPath string
	isInRepo, err := isInsideGitRepo()
	if err != nil {
		log.Printf("could not check if inside git repo: %v", err)
	}
	if isInRepo {
		repoPath, err := currentRepoRoot()
		if err != nil {
			log.Printf("could not get current repo root: %v", err)
		} else if !configHasRepo(cfg, repoPath) {
			untrackedRepoPath = repoPath
		}
	}

	dataDir, err := getDataPath()
	if err != nil {
		log.Fatalf("could not get data path: %v", err)
	}
	todayEntry, err := getTodayEntry(dataDir)
	if err != nil {
		log.Printf("could not load today entry: %v", err)
	}

	m := initialModel(*daysBack, dataDir, cfg, cfgPath, untrackedRepoPath, todayEntry)
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
