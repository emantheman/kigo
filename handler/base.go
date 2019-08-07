package handler

import (
	"fmt"
	"kigo/model"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // sets up mysql
	"github.com/gorilla/mux"
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
func New() *mux.Router {
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

	//////////////
	// HANDLERS //
	//////////////
	// A mux intelligently matches the URL of incoming reqs against registered patterns
	r := mux.NewRouter().StrictSlash(true)
	// -assets-
	r.PathPrefix("/css/").Handler(cssHandler)
	r.PathPrefix("/img/").Handler(imgHandler)
	// -root-
	r.HandleFunc("/", homeHandler)
	// -resources-
	r.HandleFunc("/haiku", getAllPoems)
	r.HandleFunc("/haiku/{id:[0-9]+}", getOnePoem)
	// -protected* resources-

	return r
}

// GOOGLE AUTH
// r.Handle("/", http.FileServer(http.Dir("template/")))
// // -google oauth2-
// r.HandleFunc("/auth/google/login", googleLogin)
// r.HandleFunc("/auth/google/callback", googleCallback)

// DUMMY DATA
// var user = model.User{Name: "nats", Email: "nats@sos.jp", Password: "1234"}
// var poem = model.Poem{Author: "kobayashi_issa", Line1: "Don’t weep, insects—", Line2: "Lovers, stars themselves,", Line3: "Must part."}
// var poem2 = model.Poem{Author: "nats", Line1: "The crow has flown away:", Line2: "Swaying in the evening sun,", Line3: "A leafless tree."}
// DB.Create(&user).Create(&poem).Create(&poem2)
