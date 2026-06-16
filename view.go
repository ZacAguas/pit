package main

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	modeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("48")).Faint(true)
	focusedPanel = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("48")).Padding(0, 1)
	dimPanel     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Faint(true)
)

func (m model) View() tea.View {
	s := ""
	switch m.view {
	case todayView:
		s += m.viewToday()
	case historyView:
		s += m.viewHistory()
	}

	v := tea.NewView(s)
	v.AltScreen = true // Fullscreen
	return v
}

func (m model) renderField(field fieldFocus, label string, t textarea.Model) string {
	style := dimPanel
	if m.focus == field {
		style = focusedPanel
	}
	return label + "\n" + style.Render(t.View())
}

// The main view when launching the app
func (m model) viewToday() string {
	s := ""

	s += m.renderField(didField, "[1] Yesterday", m.did) + "\n"
	s += m.renderField(blockedField, "[2] Blocked", m.blocked) + "\n"
	s += m.renderField(tomorrowField, "[3] Tomorrow", m.tomorrow) + "\n"

	switch m.mode {
	case normalMode:
		s += modeStyle.Render("NORMAL")
	case editMode:
		s += modeStyle.Render("EDIT")
	}
	s += " mode\n\n"

	s += m.message

	s += "\n\n" + footerStyle.Render("Move: j/k/tab | History: h | Save: s | Copy: c")
	return s
}

func (m model) viewHistory() string {
	return "history view\n\n"
}
