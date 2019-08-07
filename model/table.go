package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User is an account-holding user.
type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Salt     string
	Bio      string
}

// Poem is authored by a user.
type Poem struct {
	gorm.Model
	Author string
	User   User `gorm:"foreignkey:Author"`
	Line1  string
	Line2  string
	Line3  string
}

// FormatTime formats the poem's timestamp.
func (p *Poem) FormatTime() string {
	var (
		monthDay   string
		hourMinute = p.CreatedAt.Format("3:04PM")
		today      = time.Now()
		yesterday  = today.AddDate(0, 0, -1)
		timestamp  = p.CreatedAt
	)
	// Modifies monthDay if it is from today or yesterday
	if timestamp.Year() == today.Year() && timestamp.Month() == today.Month() && timestamp.Day() == today.Day() {
		monthDay = "Today"
	} else if timestamp.Year() == yesterday.Year() && timestamp.Month() == yesterday.Month() && timestamp.Day() == yesterday.Day() {
		monthDay = "Yesterday"
	} else {
		monthDay = p.CreatedAt.Format("Jan 2")
	}
	// Combines monthDay & hourMinute
	return monthDay + ", " + hourMinute
}

// Favorite is an association between a user and a poem that they like.
type Favorite struct {
	ID     int
	UserID int
	PoemID int
	User   User
	Poem   Poem
}
