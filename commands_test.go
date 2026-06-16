package main

import "testing"

func TestSaveEntryCmdReturnsSuccessMessage(t *testing.T) {
	dir := t.TempDir()

	e := entry{
		Date:     "2026-06-16",
		Did:      "did work",
		Blocked:  "blocked thing",
		Tomorrow: "next thing",
	}

	cmd := saveEntryCmd(dir, e)

	msg := cmd()
	got, ok := msg.(saveEntryMsg)
	if !ok {
		t.Fatalf("expected saveEntryMsg, got %T", msg)
	}
	if got.err != nil {
		t.Fatalf("expected no error, got %v", got.err)
	}

	loaded, err := loadEntry(entryFilePath(dir, e.Date))
	if err != nil {
		t.Fatal(err)
	}
	if loaded != e {
		t.Fatalf("expected %#v, got %#v", e, loaded)
	}
}
