package main

import "charm.land/bubbles/v2/key"

type todayKeyMap struct {
	Quit     key.Binding
	Edit     key.Binding
	History  key.Binding
	Save     key.Binding
	Copy     key.Binding
	Next     key.Binding
	Prev     key.Binding
	Did      key.Binding
	Blocked  key.Binding
	Tomorrow key.Binding
}

type historyKeyMap struct {
	Back key.Binding
	Open key.Binding
}

type detailKeyMap struct {
	Back key.Binding
	Copy key.Binding
	Down key.Binding
	Up   key.Binding
}

type appKeyMap struct {
	Quit key.Binding
}

var todayKeys = todayKeyMap{
	Quit:     key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "quit")),
	Edit:     key.NewBinding(key.WithKeys("i", "enter"), key.WithHelp("i/enter", "edit")),
	History:  key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "history")),
	Save:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
	Copy:     key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	Next:     key.NewBinding(key.WithKeys("j", "down", "tab"), key.WithHelp("j/↓/tab", "next")),
	Prev:     key.NewBinding(key.WithKeys("k", "up", "shift+tab"), key.WithHelp("k/↑", "prev")),
	Did:      key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "yesterday")),
	Blocked:  key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "blocked")),
	Tomorrow: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "tomorrow")),
}

var appKeys = appKeyMap{
	Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
}

var historyKeys = historyKeyMap{
	Back: key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "back")),
	Open: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
}

var detailKeys = detailKeyMap{
	Back: key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "back")),
	Copy: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	Down: key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "down")),
	Up:   key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "up")),
}

func (k todayKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.History, k.Save, k.Copy}
}

func (k todayKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Did, k.Blocked, k.Tomorrow},
		{k.Edit, k.History, k.Save, k.Copy, k.Quit},
	}
}

func (k historyKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Back}
}

func (k historyKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Open, k.Back},
	}
}

func (k detailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Down, k.Up, k.Back, k.Copy}
}

func (k detailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Down, k.Up},
		{k.Back, k.Copy},
	}
}
