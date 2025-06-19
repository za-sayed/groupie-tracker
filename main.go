package main

import (
	"groupie-tracker/handlers"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Preloaded templates map
var templates = map[string]*template.Template{}

func init() {
	loadTemplates()
}

func main() {
	// Route handlers
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// loadTemplates parses the provided template files and stores them in the templates map
func loadTemplates() {
    files, err := ioutil.ReadDir("templates")
    if err != nil {
        log.Fatalf("Error reading templates directory: %v", err)
    }
    for _, file := range files {
        if !file.IsDir() {
            tmpl, err := template.ParseFiles("templates/" + file.Name())
            if err != nil {
                log.Fatalf("Error loading template %s: %v", file.Name(), err)
            }
            templates[file.Name()] = tmpl
        }
    }
}

