package main

import "time"

func previousWorkday(today time.Time, daysBack int) time.Time {
	if daysBack < 1 {
		daysBack = 1
	}

	d := today
	for daysBack > 0 {
		d = d.AddDate(0, 0, -1)
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		daysBack--
	}
	return d
}

func commitSinceDate(today time.Time, daysBack int) string {
	return previousWorkday(today, daysBack).Format(YYYY_MM_DD)
}
