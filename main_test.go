package main

import (
	"os"
	"testing"
)

func TestGetEntryForDateReturnsExistingEntry(t *testing.T) {
	dir := t.TempDir()
	want := entry{
		Date:     "2026-06-18",
		Did:      "did work",
		Blocked:  "blocked thing",
		Tomorrow: "next thing",
	}

	if err := saveEntry(dir, want); err != nil {
		t.Fatal(err)
	}

	got, err := getEntryForDate(dir, want.Date)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("expected entry, got nil")
	}
	if *got != want {
		t.Fatalf("expected %#v, got %#v", want, *got)
	}
}

func TestGetEntryForDateReturnsNilForMissingEntry(t *testing.T) {
	got, err := getEntryForDate(t.TempDir(), "2026-06-18")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("expected nil entry, got %#v", got)
	}
}

func TestGetEntryForDateReturnsErrorForMalformedEntry(t *testing.T) {
	dir := t.TempDir()
	filePath := entryFilePath(dir, "2026-06-18")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := getEntryForDate(dir, "2026-06-18")
	if err == nil {
		t.Fatalf("expected error, got entry %#v", got)
	}
}

func TestValidateDaysBackAllowsZero(t *testing.T) {
	if err := validateDaysBack(0); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateDaysBackRejectsNegative(t *testing.T) {
	if err := validateDaysBack(-1); err == nil {
		t.Fatal("expected error, got nil")
	}
}
