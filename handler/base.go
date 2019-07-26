package handler

import (
	"fmt"
	"kigo/model"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // sets up mysql
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // sets up mysql
	"github.com/joho/godotenv"
)

// Globals
var db *gorm.DB
var err error

// Formats arguments for connection to MYSQL database.
func getConnectionArgs() string {
	// Environment vars
	dbPw := os.Getenv("DB_PASSWORD")
	dbUsr := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	// Formats connection arguments
	return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbUsr, dbPw, dbName)
}

// New returns a mux with registered patterns.
func New() *http.ServeMux {
	// Loads .env vars --NOT WORKING CORRECTLY, EXPORT REQUISITE .ENV VARS IN TERMINAL AS A TEMPORARY FIX.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Opens connection to mysql database
	db, err = gorm.Open("mysql", getConnectionArgs())
	defer db.Close()
	if err != nil {
		log.Fatal("Connection failed to open.")
	}
	log.Println("Connection established.")
	// Makes table names singular
	db.SingularTable(true)
	// Builds Tables
	db.AutoMigrate(new(model.User), new(model.Poem), new(model.Favorite))

	// A mux intelligently matches the URL of incoming reqs against registered patterns
	mux := http.NewServeMux()
	// Root
	mux.Handle("/", http.FileServer(http.Dir("template/")))
	// OauthGoogle
	mux.HandleFunc("/auth/google/login", googleLogin)
	mux.HandleFunc("/auth/google/callback", googleCallback)

	return mux
}
