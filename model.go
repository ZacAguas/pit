package main

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type viewState int // the page to show
const (
	todayView   viewState = iota // default state - showing yesterday, blocked, today panes
	historyView                  // history view, showing a list of previous entries
	detailView                   // detail view, showing rendered markdown of an entry
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

const (
	minViewWidth  = 20
	minViewHeight = 1
)

func clampMin(value, min int) int {
	if value < min {
		return min
	}
	return value
}

type model struct {
	dataDir    string
	configPath string

	config            config
	untrackedRepoPath string

	view viewState
	mode inputMode

	message string

	date string // YYYY-MM-DD

	did      textarea.Model
	blocked  textarea.Model
	tomorrow textarea.Model
	history  list.Model
	detail   viewport.Model
	help     help.Model

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

func newViewport() viewport.Model {
	v := viewport.New(viewport.WithWidth(0), viewport.WithHeight(0))
	return v
}

func initialModel(dataDir string, config config, configPath string, untrackedRepoPath string, existing *entry) model {
	did := newTextArea("What did you do yesterday?")
	blocked := newTextArea("Is anything blocking you?")
	tomorrow := newTextArea("What will you do today?")

	history := newHistoryList(nil)
	detail := newViewport()
	help := help.New()

	today := time.Now().Format(YYYY_MM_DD)
	if existing != nil {
		today = existing.Date
		did.SetValue(existing.Did)
		blocked.SetValue(existing.Blocked)
		tomorrow.SetValue(existing.Tomorrow)
	}
	return model{
		dataDir:           dataDir,
		configPath:        configPath,
		config:            config,
		untrackedRepoPath: untrackedRepoPath,

		view: todayView,
		mode: normalMode,

		message: "",

		date: today,

		did:      did,
		blocked:  blocked,
		tomorrow: tomorrow,
		history:  history,
		detail:   detail,
		help:     help,
		focus:    didField,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) resizeTextAreas() model {
	const textareaHorizontalFrame = 4 // left/right border + left/right padding

	width := m.width - textareaHorizontalFrame // could use focusedPanel.GetHorizontalFrameSize() but that couples model to styling
	width = clampMin(width, minViewWidth)
	m.did.SetWidth(width)
	m.blocked.SetWidth(width)
	m.tomorrow.SetWidth(width)
	return m
}

func (m model) resizeList() model {
	const listHorizontalFrame = 4
	const listVerticalFrame = 2

	width := m.width - listHorizontalFrame
	width = clampMin(width, minViewWidth)
	height := clampMin(m.height-listVerticalFrame, minViewHeight)

	m.history.SetSize(width, height)
	return m
}

func (m model) resizeViewport() model {
	const viewportHorizontalFrame = 4
	const footerHeight = 3

	width := clampMin(m.width-viewportHorizontalFrame, minViewWidth)
	height := clampMin(m.height-footerHeight, minViewHeight)
	m.detail.SetWidth(width)
	m.detail.SetHeight(height)

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// global msg handlers
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		m = m.resizeTextAreas()
		m = m.resizeList()
		m = m.resizeViewport()
		if m.view == detailView {
			m, _ = m.renderSelectedEntry()
		}
		return m, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, appKeys.Quit): // uncontextual exit on ctrl+c
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
	case trackRepoMsg:
		if msg.err != nil {
			m.message = "Could not track repo: " + msg.err.Error()
			return m, clearMessageAfter(3)
		}
		m.config = msg.cfg
		m.untrackedRepoPath = ""
		m.message = "Tracking repo: " + msg.repoPath
		return m, clearMessageAfter(3)
	}

	// run viewState-specific update functions
	switch m.view {
	case todayView:
		return m.updateToday(msg)
	case historyView:
		return m.updateHistory(msg)
	case detailView:
		return m.updateDetail(msg)
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
func (m model) updateTodayNormal(keyMsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, todayKeys.Quit):
		return m, tea.Quit
	case key.Matches(keyMsg, todayKeys.History):
		m.view = historyView
		return m, loadEntriesCmd(m.dataDir)
	case key.Matches(keyMsg, todayKeys.Edit):
		m.mode = editMode
		return m.applyTextAreaFocus()
	case key.Matches(keyMsg, todayKeys.Save):
		// message set/cleared in Update saveEntryMsg handler
		return m, saveEntryCmd(m.dataDir, m.currentEntry())
	case key.Matches(keyMsg, todayKeys.Copy):
		m.message = "Copied to clipboard"
		return m, tea.Batch(
			tea.SetClipboard(formatMarkdown(m.currentEntry())),
			clearMessageAfter(3),
		)
	case key.Matches(keyMsg, todayKeys.TrackRepo):
		if m.untrackedRepoPath == "" {
			return m, nil
		}
		return m, trackRepoCmd(m.configPath, m.config, m.untrackedRepoPath)
	// navigation
	case key.Matches(keyMsg, todayKeys.Next):
		m = m.focusNextField()
	case key.Matches(keyMsg, todayKeys.Prev):
		m = m.focusPrevField()
	case key.Matches(keyMsg, todayKeys.Did):
		m.focus = didField
	case key.Matches(keyMsg, todayKeys.Blocked):
		m.focus = blockedField
	case key.Matches(keyMsg, todayKeys.Tomorrow):
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
		switch {
		case key.Matches(msg, historyKeys.Back):
			m.view = todayView
			return m, nil
		case key.Matches(msg, historyKeys.Open):
			var ok bool
			m, ok = m.renderSelectedEntry()
			if ok {
				m.detail.GotoTop()
				m.view = detailView
			}
			return m, nil
		}
	}
	// pass unhandled message to list
	var cmd tea.Cmd
	m.history, cmd = m.history.Update(msg)

	return m, cmd
}

func (m model) renderSelectedEntry() (model, bool) {
	item := m.history.SelectedItem()
	e, ok := item.(entry)
	if !ok {
		return m, false
	}

	rendered, err := renderEntryMarkdown(e, m.detail.Width())
	if err != nil {
		m.message = "Could not render entry"
		return m, false
	}

	m.detail.SetContent(rendered)
	return m, true
}

func (m model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, detailKeys.Back):
			m.view = historyView
			return m, nil
		case key.Matches(msg, detailKeys.Copy):
			item := m.history.SelectedItem()
			e, ok := item.(entry)
			if !ok {
				m.message = "No entry selected"
				return m, clearMessageAfter(3)
			}
			m.message = "Copied to clipboard"
			return m, tea.Batch(
				tea.SetClipboard(formatMarkdown(e)),
				clearMessageAfter(3),
			)
		}
	}
	// pass unhandled message to viewport
	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg)
	return m, cmd
}
