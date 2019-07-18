package model

import "github.com/jinzhu/gorm"

// User is an account-holding user.
type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
	Bio      string
}

// Poem is authored by a user.
type Poem struct {
	gorm.Model
	AuthorID User
	Content  string
}

// Favorite is an association between a user and a poem that they like.
type Favorite struct {
	ID     int
	UserID User
	PoemID Poem
}
