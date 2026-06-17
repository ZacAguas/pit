package main

import (
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

type viewState int // the page to show
const (
	todayView   viewState = iota // default state - showing yesterday, blocked, today panes
	historyView                  // history view, showing a list of previous entries
)

type inputMode int // normal (navigation) or edit (text entry) - only relevant in todayView
const (
	normalMode inputMode = iota
	editMode
)

type fieldFocus int

const (
	didField fieldFocus = iota
	blockedField
	tomorrowField
)

type model struct {
	dataDir string

	view viewState
	mode inputMode

	message string

	date string // YYYY-MM-DD

	did      textarea.Model
	blocked  textarea.Model
	tomorrow textarea.Model
	history  list.Model

	focus fieldFocus

	width, height int
}

func (m model) currentEntry() entry {
	return entry{
		Date:     m.date,
		Did:      m.did.Value(),
		Blocked:  m.blocked.Value(),
		Tomorrow: m.tomorrow.Value(),
	}
}

func newTextArea(placeholder string) textarea.Model {
	t := textarea.New()

	// TODO: check these settings
	t.Placeholder = placeholder
	t.Prompt = ""
	t.ShowLineNumbers = false
	t.SetStyles(textarea.DefaultDarkStyles())
	t.SetWidth(80)
	t.SetHeight(5)

	return t
}

func newHistoryList(entries []entry) list.Model {
	l := list.New(entriesToListItems(entries), list.NewDefaultDelegate(), 0, 0)
	l.Title = "History"
	return l
}

func initialModel(dataDir string, existing *entry) model {
	did := newTextArea("What did you do yesterday?")
	blocked := newTextArea("Is anything blocking you?")
	tomorrow := newTextArea("What will you do today?")

	history := newHistoryList(nil)

	today := time.Now().Format(YYYY_MM_DD)
	if existing != nil {
		today = existing.Date
		did.SetValue(existing.Did)
		blocked.SetValue(existing.Blocked)
		tomorrow.SetValue(existing.Tomorrow)
	}
	return model{
		dataDir: dataDir,

		view: todayView,
		mode: normalMode,

		message: "",

		date: today,

		did:      did,
		blocked:  blocked,
		tomorrow: tomorrow,
		history:  history,
		focus:    didField,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) resizeTextAreas() model {
	const min = 20
	const textareaHorizontalFrame = 4 // left/right border + left/right padding

	width := m.width - textareaHorizontalFrame // could use focusedPanel.GetHorizontalFrameSize() but that couples model to styling
	if width < min {
		width = min
	}
	m.did.SetWidth(width)
	m.blocked.SetWidth(width)
	m.tomorrow.SetWidth(width)
	return m
}

func (m model) resizeList() model {
	const min = 20
	const listHorizontalFrame = 4
	width := m.width - listHorizontalFrame
	if width < min {
		width = min
	}
	height := 20

	m.history.SetSize(width, height)
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// global msg handlers
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.resizeTextAreas()
		m = m.resizeList()
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c": // uncontextual exit on ctrl+c
			return m, tea.Quit
		}
	case saveEntryMsg:
		if msg.err != nil {
			m.message = "Save failed: " + msg.err.Error()
		} else {
			m.message = "Saved"
		}
		return m, clearMessageAfter(3)
	case clearMessageMsg:
		m.message = ""
		return m, nil
	case loadEntriesMsg:
		if msg.err != nil {
			m.message = "Failed to get History"
			m.view = todayView
			return m, clearMessageAfter(3)
		}
		return m, m.history.SetItems(entriesToListItems(msg.entries))
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
	// check edit mode first, since it takes any message, not just keypress messages
	if m.mode == editMode {
		return m.updateTodayEdit(msg)
	}

	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	return m.updateTodayNormal(key)
}

func (m model) applyTextAreaFocus() (model, tea.Cmd) {
	m.did.Blur()
	m.blocked.Blur()
	m.tomorrow.Blur()

	// only apply text area focus when in edit mode
	if m.mode != editMode {
		return m, nil
	}

	switch m.focus {
	case didField:
		return m, m.did.Focus()
	case blockedField:
		return m, m.blocked.Focus()
	case tomorrowField:
		return m, m.tomorrow.Focus()
	}
	return m, nil
}

// did -> blocked -> tomorrow -> did
func (m model) focusNextField() model {
	switch m.focus {
	case didField:
		m.focus = blockedField
	case blockedField:
		m.focus = tomorrowField
	case tomorrowField:
		m.focus = didField
	}
	return m
}

// did <- blocked <- tomorrow <- did
func (m model) focusPrevField() model {
	switch m.focus {
	case didField:
		m.focus = tomorrowField
	case blockedField:
		m.focus = didField
	case tomorrowField:
		m.focus = blockedField
	}
	return m
}

// Normal mode owns navigation and app commands
func (m model) updateTodayNormal(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "q", "esc":
		return m, tea.Quit
	case "h":
		m.view = historyView
		return m, loadEntriesCmd(m.dataDir)
	case "i", "enter":
		m.mode = editMode
		return m.applyTextAreaFocus()
	case "s":
		// message set/cleared in Update saveEntryMsg handler
		return m, saveEntryCmd(m.dataDir, m.currentEntry())
	case "c":
		m.message = "Copied to clipboard"
		return m, tea.Batch(
			tea.SetClipboard(formatMarkdown(m.currentEntry())),
			clearMessageAfter(3),
		)
	// navigation
	case "j", "down", "tab":
		m = m.focusNextField()
	case "k", "up", "shift+tab":
		m = m.focusPrevField()
	case "1":
		m.focus = didField
	case "2":
		m.focus = blockedField
	case "3":
		m.focus = tomorrowField

	}
	return m, nil
}

// Reroute msg to focused text area
func (m model) updateFocusedTextArea(msg tea.Msg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.focus {
	case didField:
		m.did, cmd = m.did.Update(msg)
	case blockedField:
		m.blocked, cmd = m.blocked.Update(msg)
	case tomorrowField:
		m.tomorrow, cmd = m.tomorrow.Update(msg)
	}
	return m, cmd
}

// Edit mode only handles escape, otherwise keys go to the focused text area
func (m model) updateTodayEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	// if a keypress and "esc", change to normal mode and blur text area
	key, ok := msg.(tea.KeyPressMsg)
	if ok && key.String() == "esc" {
		m.mode = normalMode
		return m.applyTextAreaFocus()
	}

	return m.updateFocusedTextArea(msg)
}

func (m model) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q": // q goes back to today view
			m.view = todayView
			return m, nil
		}
	}
	// pass unhandled message to list
	var cmd tea.Cmd
	m.history, cmd = m.history.Update(msg)

	return m, cmd
}
