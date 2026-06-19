package main

import (
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
)

var (
	accentColor = lipgloss.Color("48")
	mutedColor  = lipgloss.Color("240")
	textColor   = lipgloss.Color("252")
	infoColor   = lipgloss.Color("104")
	warnColor   = lipgloss.Color("178")

	modeStyle         = lipgloss.NewStyle().Foreground(accentColor).Faint(true)
	focusedLabelStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	dimLabelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	focusedPanel      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accentColor).Padding(0, 1)
	dimPanel          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(mutedColor).Padding(0, 1)
	messageStyle      = lipgloss.NewStyle().Foreground(infoColor)
	warningStyle      = lipgloss.NewStyle().Foreground(warnColor)
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
	panelStyle := dimPanel
	labelStyle := dimLabelStyle
	if m.focus == field {
		panelStyle = focusedPanel
		labelStyle = focusedLabelStyle
	}
	return labelStyle.Render(label) + "\n" + panelStyle.Render(t.View())
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

	status := m.viewTodayStatus()
	if status != "" {
		s += status + "\n"
	}

	s += "\n" + m.help.View(todayKeys)
	return s
}

func (m model) viewTodayStatus() string {
	var lines []string
	if m.loadingCommits {
		lines = append(lines, messageStyle.Render("Loading commits..."))
	}
	if m.untrackedRepoPath != "" {
		lines = append(lines, warningStyle.Render("Untracked repo: "+m.untrackedRepoPath+"  [a] track"))
	}
	if m.message != "" {
		lines = append(lines, messageStyle.Render(m.message))
	}
	return strings.Join(lines, "\n")
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
