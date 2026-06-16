package main

import (
	"errors"
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
	m := initialModel(nil)

	// send an H keypress, should switch to history view
	m = update(t, m, press("h"))
	if m.view != historyView {
		t.Fatalf("expected %v, got %v", historyView, m.view)
	}
}

func TestQInHistoryReturnsToToday(t *testing.T) {
	m := initialModel(nil)
	m.view = historyView // move to history view

	// send a Q keypress, should switch to today view
	m = update(t, m, press("q"))
	if m.view != todayView {
		t.Fatalf("expected %v, got %v", todayView, m.view)
	}
}

func TestQQuitsInNormalMode(t *testing.T) {
	m := initialModel(nil)

	_, cmd := m.Update(press("q"))
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
}

func TestIEntersEditMode(t *testing.T) {
	m := initialModel(nil)

	m = update(t, m, press("i"))
	if m.mode != editMode {
		t.Fatalf("expected %v, got %v", editMode, m.mode)
	}
}

func TestEscEntersNormalMode(t *testing.T) {
	m := initialModel(nil)

	// enter edit mode
	m = update(t, m, press("i"))

	m = update(t, m, press("esc"))
	if m.mode != normalMode {
		t.Fatalf("expected %v, got %v", normalMode, m.mode)
	}
}

func TestQDoesNotQuitInEditMode(t *testing.T) {
	m := initialModel(nil)

	m.mode = editMode
	m.focus = didField
	m.did.Focus()

	next, cmd := m.Update(press("q"))
	gotModel, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	// no command - no quit was issued
	if cmd != nil {
		if _, isQuit := cmd().(tea.QuitMsg); isQuit {
			t.Fatal("expected 'q' not to quit but a tea.Quit command was returned")
		}
	}

	if got := gotModel.did.Value(); got != "q" {
		t.Fatalf("expected q to be inserted into 'did' field, got %q", got)
	}
}

func TestNextFieldNavigationCycles(t *testing.T) {
	m := initialModel(nil)

	m = update(t, m, press("j"))
	if m.focus != blockedField {
		t.Fatalf("expected blockedField focused, got %v", m.focus)
	}

	m = update(t, m, press("j"))
	if m.focus != tomorrowField {
		t.Fatalf("expected tomorrowField focused, got %v", m.focus)
	}

	m = update(t, m, press("j"))
	if m.focus != didField {
		t.Fatalf("expected didField focused, got %v", m.focus)
	}
}

func TestPreviousFieldNavigationCycles(t *testing.T) {
	m := initialModel(nil)

	m = update(t, m, press("k"))
	if m.focus != tomorrowField {
		t.Fatalf("expected tomorrowField focused, got %v", m.focus)
	}

	m = update(t, m, press("k"))
	if m.focus != blockedField {
		t.Fatalf("expected blockedField focused, got %v", m.focus)
	}

	m = update(t, m, press("k"))
	if m.focus != didField {
		t.Fatalf("expected didField focused, got %v", m.focus)
	}
}

func TestNumberKeysJumpToFields(t *testing.T) {
	m := initialModel(nil)

	m = update(t, m, press("2"))
	if m.focus != blockedField {
		t.Fatalf("expected blockedField focused, got %v", m.focus)
	}

	m = update(t, m, press("3"))
	if m.focus != tomorrowField {
		t.Fatalf("expected tomorrowField focused, got %v", m.focus)
	}

	m = update(t, m, press("1"))
	if m.focus != didField {
		t.Fatalf("expected didField focused, got %v", m.focus)
	}
}

func TestCurrentEntryUsesTextFields(t *testing.T) {
	m := initialModel(nil)

	const did = "did work"
	const blocked = "blocked thing"
	const tomorrow = "next thing"
	m.did.SetValue(did)
	m.blocked.SetValue(blocked)
	m.tomorrow.SetValue(tomorrow)

	got := m.currentEntry()
	if got.Did != did {
		t.Fatalf("expected %q got %q", did, got.Did)
	}
	if got.Blocked != blocked {
		t.Fatalf("expected %q got %q", blocked, got.Blocked)
	}
	if got.Tomorrow != tomorrow {
		t.Fatalf("expected %q got %q", tomorrow, got.Tomorrow)
	}
}

func TestCCopiesCurrentEntryShowsMessage(t *testing.T) {
	m := initialModel(nil)

	next, cmd := m.Update(press("c"))
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	if got.message != "Copied to clipboard" {
		t.Fatalf("expected copy message, got %q", got.message)
	}

	if cmd == nil {
		t.Fatal("expected command, got nil")
	}
}

func TestClearMessageClearsMessage(t *testing.T) {
	m := initialModel(nil)
	m.message = "Copied to clipboard"

	next, _ := m.Update(clearMessageMsg{})
	got, ok := next.(model)
	if !ok {
		t.Fatalf("expected model, got %T", next)
	}

	if got.message != "" {
		t.Fatalf("expected message to be cleared, got %q", got.message)
	}
}

func TestSaveEntrySuccessShowsSavedMessage(t *testing.T) {
	m := initialModel(nil)

	next, _ := m.Update(saveEntryMsg{err: nil})
	got := next.(model)

	if got.message != "Saved" {
		t.Fatalf("expected Saved message, got %q", got.message)
	}
}

func TestSaveEntryErrorShowsFailureMessage(t *testing.T) {
	m := initialModel(nil)

	next, _ := m.Update(saveEntryMsg{err: errors.New("permission denied")})
	got := next.(model)

	want := "Save failed: permission denied"
	if got.message != want {
		t.Fatalf("expected %q, got %q", want, got.message)
	}
}
