package main

import (
	"log"
	"net/http"

	"kigo/handler"
	"kigo/session"
	_ "kigo/session/provider/memory" // registers memory provider
)

func init() {
	// Initializes a global session manager
	session.MegaManager, err = session.NewManager("memory", "oauthstate", 3600)
	if err != nil {
		log.Printf("Error creating global session-manager.")
	}
	// Starts GC cycle
	go session.MegaManager.GC()
}

func main() {
	// Creates a simple http server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.New(),
	}

	// Runs server
	log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}
}

// // Global vars
// var db *gorm.DB
// var err error

// // Load templates
// // or use: template.ParseFiles("tmpl/edit.html", "tmpl/view.html") and list view paths as args
// var templates = template.Must(template.ParseGlob("tmpl/*"))

// // Renders an html template.
// func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
// 	// Sends a templated html response to the ResponseWriter
// 	err := templates.ExecuteTemplate(w, tmpl+".html", data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// // Redirects to "/view/home"
// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "home", nil)
// }

// // LoadPoem retrieves a poem by id.
// func loadPoem(id int) (Poem, error) {
// 	var poem Poem
// 	// SELECT * FROM poem WHERE id = <id>
// 	err = db.Model(&Poem{}).Where("id = ?", id).Take(&poem).Error
// 	return poem, err
// }

// // ValidPoem validates whether poem adheres to haiku constraints.
// func validPoem(lines []string) bool {
// 	return true
// }

// // CountFavs counts the number of favorites a poem has.
// func countFavs(id int) (int, error) {
// 	var count int
// 	// SELECT count(*) FROM poem WHERE id = <id>
// 	err = db.Model(&Poem{}).Where("id = ?", id).Count(&count).Error
// 	return count, err
// }

// // Handles viewing a poem.
// func viewHandler(w http.ResponseWriter, r *http.Request, id ...int) {
// 	// Handles no ID being passed
// 	if len(id) == 0 {
// 		http.Error(w, "A poem with this ID doesn't exist.", http.StatusBadRequest)
// 		return
// 	}
// 	// Loads poem by id
// 	p, err := loadPoem(id[0])
// 	// Redirect to new Poem if poem doesn't yet exist
// 	if err != nil {
// 		http.Redirect(w, r, "/post/", http.StatusFound)
// 		return
// 	}
// 	// Sends a templated response to the writer
// 	renderTemplate(w, "view", p)
// }

// // Handles editing a wiki.
// func editHandler(w http.ResponseWriter, r *http.Request, id ...int) {
// 	var poem Poem
// 	// If an ID is passed, edits a poem
// 	if len(id) != 0 {
// 		// Loads page by title
// 		poem, err = loadPoem(id[0])
// 		if err != nil {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		// If ID is nil, creates a poem
// 	} else {
// 		// Author's id as a string -- TEMPORARY SOLUTION TO USER_ID STORAGE PROBLEM
// 		strID := r.FormValue("author-id")
// 		// Converts strID to int
// 		authorID, err := strconv.Atoi(strID)
// 		if err != nil {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		// Creates poem
// 		poem = Poem{AuthorID: authorID, Line1: r.FormValue("line-1"), Line2: r.FormValue("line-2"), Line3: r.FormValue("line-3")}
// 		// Saves poem to database
// 		err = db.Create(&poem).Error
// 		if err != nil {
// 			http.Error(w, "Error on CreatePoem.", http.StatusBadRequest)
// 			return
// 		}
// 		// Redirects to poem-view
// 		http.Redirect(w, r, "/view/"+strID, http.StatusFound)
// 		return
// 	}
// 	// Sends a templated response to the writer
// 	renderTemplate(w, "edit", poem)
// }

// // Handles saving a poem.
// func saveHandler(w http.ResponseWriter, r *http.Request, id ...int) {
// 	poem := r.FormValue("poem") // get poem
// 	// ---CONTENT VALIDATION REQUIRED---
// 	// Update Poem

// 	// Go to view
// 	http.Redirect(w, r, "/view/"+strconv.Itoa(id[0]), http.StatusFound)
// }

// // Valid url
// var validPath = regexp.MustCompile("^/(edit|save|view|poet)/([a-zA-Z0-9]+)$")

// // Retrieves and error-checks title, and returns a HandlerFunc.
// func makeHandler(fn func(http.ResponseWriter, *http.Request, ...int)) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Tests path against regex
// 		m := validPath.FindStringSubmatch(r.URL.Path)
// 		if m == nil {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		// Displays path
// 		fmt.Println(m)
// 		// Converts id to integer
// 		id, err := strconv.Atoi(m[2])
// 		if err != nil {
// 			http.Error(w, "Unable to convert ID to string.", http.StatusBadRequest)
// 		}
// 		// Calls handler with id
// 		fn(w, r, id)
// 	}
// }

// // Formats arguments for connection to MYSQL database.
// func getConnectionArgs() string {
// 	// Loads environment vars
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env variables.")
// 	}
// 	// Environment variables
// 	dbPw := os.Getenv("DB_PASSWORD")
// 	dbUsr := os.Getenv("DB_USER")
// 	dbName := os.Getenv("DB_NAME")
// 	// Formats connection arguments
// 	return fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbUsr, dbPw, dbName)
// }

// func main() {
// 	// Opens connection to mysql database
// 	db, err = gorm.Open("mysql", getConnectionArgs())
// 	defer db.Close()
// 	if err != nil {
// 		log.Fatal("Connection failed to open.")
// 	}
// 	log.Println("Connection established.")

// 	// Makes table names singular
// 	db.SingularTable(true)

// 	// Builds Tables
// 	db.AutoMigrate(new(User), new(Poem), new(Favorite))

// 	// Test create user
// 	// var u = User{Name: "tboy", Password: "test", Email: "t@tboy", Bio: "I'm a teemster"}
// 	var haiku = Poem{AuthorID: 1, Line1: "Growing at home", Line2: "A dragon", Line3: "No one sees"}
// 	db.Create(&haiku)

// 	// Serves everything in the css and img folder as a file
// 	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
// 	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

// 	// Tells http-package to handle all reqs to <path> with <handler>
// 	http.HandleFunc("/", homeHandler)
// 	http.HandleFunc("/view/", makeHandler(viewHandler))
// 	http.HandleFunc("/edit/", makeHandler(editHandler))
// 	http.HandleFunc("/save/", makeHandler(saveHandler))

// 	// L&S listens on port :8080; nil is placeholder for a middleware
// 	log.Fatal(http.ListenAndServe(":8080", nil)) // wrap w/ log.Fatal in case of error
// }
