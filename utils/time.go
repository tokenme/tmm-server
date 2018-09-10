package utils

import "time"

func TimeToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func YearRemainSeconds() float64 {
	now := time.Now()
	thisYear := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	nextYear := thisYear.AddDate(1, 0, 0)
	du := nextYear.Sub(now)
	return du.Seconds()
}
