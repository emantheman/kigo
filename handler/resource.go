package handler

import (
	"fmt"
	"html/template"
	"kigo/model"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Data object for template rendering
type Data struct {
	IsAuthenticated bool // whether the user is authenticated
	Username        string
	Poems           []model.Poem // slice of poems
}

// or use: template.ParseFiles("tmpl/edit.html", "tmpl/view.html") and list view paths as args
var (
	// Load templates
	templates = template.Must(template.ParseGlob("tmpl/*"))
	// Valid url
	validPath = regexp.MustCompile("^/(haiku|author)/([a-zA-Z0-9]+)$")
)

// Renders an html template.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Sends a templated html response to the ResponseWriter
	err := templates.ExecuteTemplate(w, tmpl+".html", &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Writes success msg to console
	log.Printf("Successfully rendered template '%s.'", tmpl)
}

// Renders * poems to client.
func getAllPoems(w http.ResponseWriter, r *http.Request) {
	var poems = []model.Poem{} // slice of poems
	// Retrieves * poems from DB
	DB.Order("created_at DESC").Find(&poems)
	// Creates data object for rendering
	// data := Data{Poems: poems}
	// Renders poems to client
	renderTemplate(w, "allpoems", poems)
}

// Wraps functions that receive IDs; returns an http.HandlerFunc that will validate the URL according to a preset regex, parse it for an ID, then pass the writer, request, and ID to a custom handler.
func makeDynamicHandlerFunc(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Tests path against regex
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		// Displays path
		fmt.Println(m)
		// Converts id to integer
		id, err := strconv.Atoi(m[2])
		if err != nil {
			http.Error(w, "Unable to convert ID to string.", http.StatusBadRequest)
		}
		// Calls handler with id
		fn(w, r, id)
	}
}

func getOnePoem(w http.ResponseWriter, r *http.Request, id int) {
	var poem = model.Poem{}
}
