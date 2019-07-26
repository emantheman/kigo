package handler

import (
	"kigo/model"
	"log"
	"net/http"
	"text/template"
)

// Data object for template rendering
type Data struct {
	IsAuthenticated bool // whether the user is authenticated
	Username        string
	Poems           []*model.Poem // slice of poems
}

// Load templates
// or use: template.ParseFiles("tmpl/edit.html", "tmpl/view.html") and list view paths as args
var templates = template.Must(template.ParseGlob("tmpl/*"))

// Renders an html template.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Sends a templated html response to the ResponseWriter
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Writes success msg to console
	log.Printf("Successfully rendered template %s.", tmpl)
}

// Renders * poems to client.
func getAllPoems(w http.ResponseWriter, r *http.Request) {
	var poems = []*model.Poem{} // slice of poem ptrs
	// Retrieves * poems from DB
	DB.Find(&poems)
	// Creates data object for rendering
	data := &Data{Poems: poems}
	// Renders poems to client
	renderTemplate(w, "allpoems", data)
}
