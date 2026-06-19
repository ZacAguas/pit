package main

import (
	"strconv"
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
	editColor   = lipgloss.Color("42")

	normalModeStyle   = lipgloss.NewStyle().Foreground(infoColor).Bold(true)
	editModeStyle     = lipgloss.NewStyle().Foreground(editColor).Bold(true)
	focusedLabelStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	dimLabelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	focusedPanel      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accentColor).Padding(0, 1)
	dimPanel          = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(mutedColor).Padding(0, 1)
	messageStyle      = lipgloss.NewStyle().Foreground(infoColor)
	warningStyle      = lipgloss.NewStyle().Foreground(warnColor)
	lineCountStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
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
	panel := panelStyle.Render(t.View())
	counter := lineCountStyle.Render(textAreaLineCounter(t, m.focus == field))
	return labelStyle.Render(label) + "\n" + panel + "\n" + lipgloss.PlaceHorizontal(lipgloss.Width(panel), lipgloss.Right, counter)
}

func textAreaLineCounter(t textarea.Model, selected bool) string {
	lineCount := t.LineCount()
	if !selected {
		if lineCount == 1 {
			return "1 line"
		}
		return strconv.Itoa(lineCount) + " lines"
	}
	currentLine := t.Line() + 1
	return "line " + strconv.Itoa(currentLine) + "/" + strconv.Itoa(lineCount)
}

// The main view when launching the app
func (m model) viewToday() string {
	s := ""

	s += m.renderField(didField, "[1] Yesterday", m.did) + "\n"
	s += m.renderField(blockedField, "[2] Blocked", m.blocked) + "\n"
	s += m.renderField(tomorrowField, "[3] Tomorrow", m.tomorrow) + "\n"

	switch m.mode {
	case normalMode:
		s += normalModeStyle.Render("NORMAL")
	case editMode:
		s += editModeStyle.Render("EDIT  ")
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
