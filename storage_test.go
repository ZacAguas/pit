package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadEntry(t *testing.T) {
	const date = "2020-12-25"
	dir := t.TempDir()
	e := entry{
		Date:     "2020-12-25",
		Did:      "a",
		Blocked:  "b",
		Tomorrow: "c",
	}

	err := saveEntry(dir, e)
	if err != nil {
		t.Fatal(err)
	}

	got, err := loadEntry(entryFilePath(dir, e.Date))
	if err != nil {
		t.Fatal(err)
	}

	if e != got {
		t.Fatalf("expected %#v, got %#v", e, got)
	}
}

func TestGetAllEntriesReturnsEntriesSortedNewestFirst(t *testing.T) {
	dir := t.TempDir()

	entries := []entry{
		{Date: "2020-12-24", Did: "older"},
		{Date: "2020-12-26", Did: "newer"},
		{Date: "2020-12-25", Did: "middle"},
	}

	for _, e := range entries {
		if err := saveEntry(dir, e); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("ignore me"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := getAllEntries(dir)
	if err != nil {
		t.Fatal(err)
	}

	wantDates := []string{"2020-12-26", "2020-12-25", "2020-12-24"}
	if len(got) != len(wantDates) {
		t.Fatalf("expected %d entries, got %d: %#v", len(wantDates), len(got), got)
	}

	for i, want := range wantDates {
		if got[i].Date != want {
			t.Fatalf("entry %d: expected date %q, got %q", i, want, got[i].Date)
		}
	}
}
