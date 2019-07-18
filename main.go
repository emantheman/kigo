package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

// Page is a representation of a wiki page.
type Page struct {
	Title string
	Body  []byte
}

// Saves the Page to a text file.
func (p *Page) save() error {
	filename := p.Title + ".txt" // file is named after title
	// Writes p.Body to a file with title p.Title
	return ioutil.WriteFile("data/"+filename, p.Body, 0600) // code 0600 restricts read-write permissions to the current user
}

// Loads the Page.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	// Read the text file
	body, err := ioutil.ReadFile("data/" + filename)
	// If error occurs, let caller deal with it
	if err != nil {
		return nil, err
	}
	// Send Page
	return &Page{Title: title, Body: body}, nil
}

// Load templates
// or use: template.ParseFiles("tmpl/edit.html", "tmpl/view.html") and list view paths as args
var templates = template.Must(template.ParseGlob("tmpl/*"))

// Renders an html template.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// Sends a templated html response to the ResponseWriter
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Redirects to "/view/home"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", nil)
}

// Handles viewing a wiki.
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	// Loads page by title
	p, err := loadPage(title)
	// Redirect to new Page if page doesn't yet exist
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	// Sends a templated response to the writer
	renderTemplate(w, "view", p)
}

// Handles editing a wiki.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	// Loads page by title
	p, err := loadPage(title)
	// Creates empty Page if page does not yet exist
	if err != nil {
		p = &Page{Title: title}
	}
	// Sends a templated response to the writer
	renderTemplate(w, "edit", p)
}

// Handles saving a wiki.
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	// Get body
	body := r.FormValue("body")
	// Make Page
	p := &Page{Title: title, Body: []byte(body)}
	// Save Page
	if err := p.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Go to view
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Valid filename
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// Retrieves and error-checks title, and returns a HandlerFunc.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Println(m)
		fn(w, r, m[2])
	}
}

// Retrieves and formats environment variables for connection to MYSQL database.
func getConnectionArgs() string {
	// Loads environment vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env variables.")
	}
	// Environment variables
	dbPw := os.Getenv("DB_PASSWORD")
	dbUsr := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	// Formats connection arguments
	return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbUsr, dbPw, dbName)
}

func main() {
	// Opens connection to mysql database
	db, err := gorm.Open("mysql", getConnectionArgs())
	defer db.Close()
	if err != nil {
		log.Fatal("Connection failed to open.")
	}
	log.Println("Connection established.")

	// Serves everything in the css and img folder as a file
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	// Tells http-package to handle all reqs to <path> with <handler>
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	// L&S listens on port :8080; nil is placeholder for a middleware
	log.Fatal(http.ListenAndServe(":8080", nil)) // wrap w/ log.Fatal in case of error
}
