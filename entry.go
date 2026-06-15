package main

type entry struct {
	Date     string `json:"date"` // Stored as ISO — "YYYY-MM-DD"
	Did      string `json:"did"`
	Blocked  string `json:"blocked"`
	Tomorrow string `json:"tomorrow"`
}
