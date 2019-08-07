package handler

import "time"

func formatTime(t time.Time) string {
	var (
		monthDay   string
		hourMinute = t.Format("3:04PM")
		today      = time.Now()
		yesterday  = today.AddDate(0, 0, -1)
		timestamp  = t
	)
	// Modifies monthDay if it is from today or yesterday
	if timestamp.Year() == today.Year() && timestamp.Month() == today.Month() && timestamp.Day() == today.Day() {
		monthDay = "Today"
	} else if timestamp.Year() == yesterday.Year() && timestamp.Month() == yesterday.Month() && timestamp.Day() == yesterday.Day() {
		monthDay = "Yesterday"
	} else {
		monthDay = t.Format("Jan 2")
	}
	// Combines monthDay & hourMinute
	return monthDay + ", " + hourMinute
}
