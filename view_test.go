package main

import "testing"

func TestTextAreaLineCounterShowsCurrentLineWhenSelected(t *testing.T) {
	area := newTextArea("")
	area.SetValue("one\ntwo")

	got := textAreaLineCounter(area, true)
	want := "line 2/2"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestTextAreaLineCounterShowsLineCountWhenNotSelected(t *testing.T) {
	area := newTextArea("")
	area.SetValue("one\ntwo")

	got := textAreaLineCounter(area, false)
	want := "2 lines"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
