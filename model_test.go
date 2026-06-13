package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func press(key string) tea.KeyPressMsg {
	if len(key) == 1 {
		r := rune(key[0])
		return tea.KeyPressMsg(tea.Key{
			Text: key,
			Code: r,
		})
	}
	switch key {
	case "esc":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEsc})
	case "enter":
		return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	default:
		return tea.KeyPressMsg(tea.Key{Text: key})
	}
}

func update(t *testing.T, m model, msg tea.Msg) model {
	t.Helper() // mark function as a test helper - skipped when printing line/file info

	next, _ := m.Update(msg)
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	return got
}

func TestHOpensHistory(t *testing.T) {
	m := newModel()

	// send an H keypress, should switch to history view
	m = update(t, m, press("h"))
	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
}

func TestQInHistoryReturnsToToday(t *testing.T) {
	m := newModel()
	m.view = historyView // move to history view

	// send a Q keypress, should switch to today view
	m = update(t, m, press("q"))
	if m.view != todayView {
		t.Fatalf("expected %v, got %v", todayView, m.view)
	}
}

func TestQQuitsInNormalMode(t *testing.T) {
	m := newModel()

	_, cmd := m.Update(press("q"))
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}
