package handler

import (
	"html/template"
	"kigo/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Data object for template rendering
type Data struct {
	IsAuthenticated bool // whether the user is authenticated
	Username        string
	Poems           []model.Poem // slice of poems
}

var (
	// Functions for use within templates
	funcmap = template.FuncMap{
		"safe":       safe,
		"attr":       attr,
		"formattime": formatTime,
		"randbgc":    backgroundColor,
	}
	// Templates from tmpl/ directory
	templates = template.Must(template.New("tmp").Funcs(funcmap).ParseGlob("tmpl/*"))
	// Asset handlers
	cssHandler = http.StripPrefix("/css/", http.FileServer(http.Dir("css")))
	imgHandler = http.StripPrefix("/img/", http.FileServer(http.Dir("img")))
)

// Renders an html template.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Sends a templated html response to the ResponseWriter
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Writes success msg to console
	log.Printf("Successfully rendered template '%s.'", tmpl)
}

// Renders all poems to client.
func getAllPoems(w http.ResponseWriter, r *http.Request) {
	var poems []model.Poem // slice of poems
	// Retrieves * poems from DB
	DB.Order("created_at DESC").Find(&poems)
	// Renders poems to client
	renderTemplate(w, "poemIndex", poems)
}

// Renders one poem to client.
func getOnePoem(w http.ResponseWriter, r *http.Request) {
	var poem model.Poem
	// Retrieves ID from url
	id := mux.Vars(r)["id"]
	// Retrieves poem by ID
	DB.Where("id = ?", id).First(&poem)
	// Renders poem to client
	renderTemplate(w, "poemShow", poem)
}
