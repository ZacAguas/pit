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
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "h":
			m.view = historyView
		case "q", "esc":
			if m.mode == normalMode {
				return m, tea.Quit // only accept q/esc as quit in normal mode
			}
		}
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
