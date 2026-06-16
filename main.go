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
		getDataPath()
	}

	todayEntry, err := getTodayEntry()
	if err != nil {
		log.Printf("could not load today entry: %v", err)
	}
	m := initialModel(todayEntry)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}

func getTodayEntry() (*entry, error) {
	today := time.Now().Format(YYYY_MM_DD)
	filePath := entryFilePath(getDataPath(), today)
	e, err := loadEntry(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}
