package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	m := newModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}
