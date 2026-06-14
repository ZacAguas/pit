package main

import (
	tea "charm.land/bubbletea/v2"
)

type viewState uint // the page to show
const (
	todayView   = viewState(iota) // default state - showing yesterday, blocked, today panes
	historyView                   // history view, showing a list of previous entries
)

type inputMode uint // normal (navigation) or edit (text entry) - only relevant in todayView
const (
	normalMode = inputMode(iota)
	editMode
)

type model struct {
	view viewState
	mode inputMode
}

func newModel() model {
	return model{
		view: todayView,
		mode: normalMode,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// global msg handlers
	// TODO: window resize handler
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c": // uncontextual exit on ctrl+c
			return m, tea.Quit
		}
	}

	// run viewState-specific update functions
	switch m.view {
	case todayView:
		return m.updateToday(msg)
	case historyView:
		return m.updateHistory(msg)
	}
	return m, nil
}

func (m model) updateToday(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch m.mode {
	case normalMode:
		return m.updateTodayNormal(key)
	case editMode:
		return m.updateTodayEdit(key)
	}

	return m, nil
}

// Normal mode owns navigation and app commands
func (m model) updateTodayNormal(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q", "esc":
		return m, tea.Quit
	case "h":
		m.view = historyView
	case "i", "enter":
		m.mode = editMode
	}
	return m, nil
}

// Edit mode only handles escape, otherwise keys go to the focused text area
func (m model) updateTodayEdit(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc": // no 'q' as to not swallow text keypresses when editing
		m.mode = normalMode
	}
	return m, nil
}

func (m model) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc": // q/esc goes back to today view
			m.view = todayView
		}
	}
	return m, nil
}
