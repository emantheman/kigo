package model

import (
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

// Favorite is an association between a user and a poem that they like.
type Favorite struct {
	ID     int
	UserID int
	PoemID int
	User   User
	Poem   Poem
}
