package handler

import (
	"html/template"
	"math/rand"
	"time"
)

// Formats a timestamp.
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

// Tip: If attr() or safe() are not used on unsafe attirbutes and content, the special value "ZgotmplZ" will replace them.

// Turns a string into an html attribute.
func attr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

// Turns a string into an html tag.
func safe(s string) template.HTML {
	return template.HTML(s)
}

// Generates a random bgc.
func backgroundColor() string {
	// List of hex background-colors
	var bgcolors = []string{"#ff9fd3", "#9386e6", "#647ce4", "#4aac69", "#654c7f", "#8568a4", "#687ca4", "#71a610", "#ae915b", "#49c16f", "#ae6b95", "#c38fb0", "#b8004a", "#8f3a5c", "#bf4848", "#219091", "#6d835a", "#ac7b8c", "#3d004f"}
	// Initializes local pseudorandom generator
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// Gets a random color from bgcolors
	return bgcolors[r.Intn(len(bgcolors))]
}
