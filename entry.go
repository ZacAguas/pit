package main

import (
	"strings"

	"charm.land/bubbles/v2/list"
)

type entry struct {
	Date     string `json:"date"` // Stored as ISO — "YYYY-MM-DD"
	Did      string `json:"did"`
	Blocked  string `json:"blocked"`
	Tomorrow string `json:"tomorrow"`
}

// Implement List.Item interface for history list
func (e entry) FilterValue() string {
	return e.Date + " " + e.Did + " " + e.Blocked + " " + e.Tomorrow
}

func (e entry) Title() string {
	return formatDateForHeading(e.Date)
}

func (e entry) Description() string {
	for _, value := range []string{e.Did, e.Blocked, e.Tomorrow} {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		return strings.Split(value, "\n")[0]
	}
	return "No details"
}

func entriesToListItems(entries []entry) []list.Item {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = e
	}
	return items
}
