package main

import (
	"strings"
	"testing"
)

func TestMarkdownOutput(t *testing.T) {
	entry := entry{
		Date:     "2020-12-25",
		Did:      "a",
		Blocked:  "b",
		Tomorrow: "c",
	}

	output := formatMarkdown(entry)
	expected := "### Standup — December 25, 2020\n\n" +
		"**Yesterday**\n" +
		"a\n\n" +
		"**Blocked**\n" +
		"b\n\n" +
		"**Today**\n" +
		"c\n"

	if strings.Compare(output, expected) != 0 {
		t.Fatalf("expected %q got %q", expected, output)
	}
}
