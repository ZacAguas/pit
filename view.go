package main

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
)

var (
	modeStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("48")).Faint(true)
	focusedPanel = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("48")).Padding(0, 1)
	dimPanel     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
	messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("104"))
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("178"))
)

func (m model) View() tea.View {
	s := ""
	switch m.view {
	case todayView:
		s += m.viewToday()
	case historyView:
		s += m.viewHistory()
	case detailView:
		s += m.viewDetail()
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

	if m.loadingCommits {
		s += messageStyle.Render("Loading commits...") + "\n"
	}

	if m.message != "" {
		s += messageStyle.Render(m.message) + "\n"
	}

	if m.untrackedRepoPath != "" {
		s += warningStyle.Render("Untracked repo: " + m.untrackedRepoPath)
	}

	s += "\n\n" + m.help.View(todayKeys)
	return s
}

func (m model) viewHistory() string {
	s := ""
	s += m.history.View()
	if m.message != "" {
		s += "\n\n" + messageStyle.Render(m.message)
	}
	s += "\n\n" + m.help.View(historyKeys)
	return s
}

func (m model) viewDetail() string {
	message := ""
	if m.message != "" {
		message = "\n\n" + messageStyle.Render(m.message)
	}

	return m.detail.View() + message + "\n\n" + m.help.View(detailKeys)
}

func renderEntryMarkdown(e entry, width int) (string, error) {
	if width <= 0 {
		width = 80
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	return renderer.Render(formatMarkdown(e))
}
