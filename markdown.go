package main

import (
	"log"
	"strings"
	"time"
)

func formatMarkdown(e entry) string {
	var b strings.Builder

	b.WriteString("# Standup")
	if e.Date != "" {
		b.WriteString(" — ")
		b.WriteString(formatDateForHeading(e.Date))
	}
	b.WriteString("\n\n")

	b.WriteString("**Yesterday**\n")
	b.WriteString(e.Did)
	b.WriteString("\n\n")

	b.WriteString("**Blocked**\n")
	b.WriteString(e.Blocked)
	b.WriteString("\n\n")

	b.WriteString("**Today**\n")
	b.WriteString(e.Tomorrow)
	b.WriteString("\n")

	return b.String()
}

// Assuming YYYY-MM-DD
func formatDateForHeading(date string) string {
	parsedTime, err := time.Parse(YYYY_MM_DD, date)
	if err != nil {
		log.Printf("error parsing time: %v", err)
		return date
	}

	formattedDate := parsedTime.Format("January 2, 2006")
	return formattedDate
}

func joinSections(sections ...string) string {
	var nonEmpty []string
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section != "" {
			nonEmpty = append(nonEmpty, section)
		}
	}
	return strings.Join(nonEmpty, "\n\n")
}
