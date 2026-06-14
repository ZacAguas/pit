package main

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) View() tea.View {
	s := "pit\n\n"
	switch m.view {
	case todayView:
		s += m.viewToday()
	case historyView:
		s += m.viewHistory()
	}

	v := tea.NewView(s)
	// v.AltScreen = true // Fullscreen
	return v
}

// The main view when launching the app
func (m model) viewToday() string {
	s := ""
	switch m.mode {
	case normalMode:
		s += "NORMAL"
	case editMode:
		s += "EDIT"
	}
	s += " mode\n"

	return s + "pane view\n\n"
}

func (m model) viewHistory() string {
	return "history view\n\n"
}
