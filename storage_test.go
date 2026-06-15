package main

import (
	"reflect"
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

	if !reflect.DeepEqual(e, got) {
		t.Fatalf("expected %#v, got %#v", e, got)
	}
}
