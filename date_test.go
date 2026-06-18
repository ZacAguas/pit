package main

import (
	"testing"
	"time"
)

func TestPreviousWorkday(t *testing.T) {
	tests := []struct {
		name     string
		today    string
		daysBack int
		want     string
	}{
		{
			name:     "Tuesday returns Monday",
			today:    "2026-06-16",
			daysBack: 1,
			want:     "2026-06-15",
		},
		{
			name:     "Monday returns Friday",
			today:    "2026-06-22",
			daysBack: 1,
			want:     "2026-06-19",
		},
		{
			name:     "Wednesday two days back returns Monday",
			today:    "2026-06-17",
			daysBack: 2,
			want:     "2026-06-15",
		},
		{
			name:     "Monday two days back returns Thursday",
			today:    "2026-06-22",
			daysBack: 2,
			want:     "2026-06-18",
		},
		{
			name:     "Zero days back defaults to one",
			today:    "2026-06-16",
			daysBack: 0,
			want:     "2026-06-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			today := mustParseDate(t, tt.today)
			want := mustParseDate(t, tt.want)

			got := previousWorkday(today, tt.daysBack)
			if !got.Equal(want) {
				t.Fatalf("expected %v, got %v", want, got)
			}
		})
	}
}

func mustParseDate(t *testing.T, value string) time.Time {
	t.Helper()

	date, err := time.Parse(YYYY_MM_DD, value)
	if err != nil {
		t.Fatal(err)
	}
	return date
}
