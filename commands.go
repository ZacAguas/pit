package main

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type saveEntryMsg struct {
	err error
}

func saveEntryCmd(dir string, e entry) tea.Cmd {
	return func() tea.Msg {
		return saveEntryMsg{
			err: saveEntry(dir, e),
		}
	}
}

type clearMessageMsg struct{}

func clearMessageAfter(seconds int) tea.Cmd {
	return tea.Tick(time.Duration(seconds)*time.Second, func(time.Time) tea.Msg {
		return clearMessageMsg{}
	})
}
