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

// DB is a gorm database
var DB *gorm.DB
var err error

// Redirects to "/home"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", nil)
}

// Formats arguments for connection to MYSQL database.
func getConnectionArgs() string {
	// Environment vars
	DBPw := os.Getenv("DB_PASSWORD")
	DBUsr := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_NAME")
	// Formats connection arguments
	return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", DBUsr, DBPw, DBName)
}

// New registers and returns a mux.
func New() *http.ServeMux {
	//////////
	// .ENV //
	//////////
	// Loads .env vars
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//////////////
	// DATABASE //
	//////////////
	// Opens connection to mysql database
	DB, err = gorm.Open("mysql", getConnectionArgs())
	if err != nil {
		log.Fatal("Connection failed to open.")
	}
	log.Println("Connection established.")
	// Builds Tables
	DB.AutoMigrate(new(model.User), new(model.Poem), new(model.Favorite))
	// DUMMY DATA
	// var poem = model.Poem{AuthorID: 1, Line1: "The dog licks my feet", Line2: "It ignores my loud commands", Line3: "With the happiest smile"}
	// var poem2 = model.Poem{AuthorID: 1, Line1: "In the blue-green pond", Line2: "Two girls swim nude", Line3: "Their skin is white!"}
	// DB.Create(&poem).Create(&poem2)

	//////////////
	// HANDLERS //
	//////////////
	// A mux intelligently matches the URL of incoming reqs against registered patterns
	mux := http.NewServeMux()
	// -css/img-
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	// -root-
	mux.Handle("/", http.FileServer(http.Dir("template/")))
	mux.HandleFunc("/home", homeHandler)
	// -google oauth2-
	mux.HandleFunc("/auth/google/login", googleLogin)
	mux.HandleFunc("/auth/google/callback", googleCallback)
	// -resources
	mux.HandleFunc("/haiku", getAllPoems)
	// -protected* resources-

	return mux
}
