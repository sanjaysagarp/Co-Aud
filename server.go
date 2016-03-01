package main

import (
	"html/template"
	"net/http"
	"fmt"
	"os"
	"path"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("./app/views/header.html", "./app/views/footer.html", "./app/views/main.html"))

//A Page structure
type Page struct {
	Title string
}

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

//The handlers.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "main", &Page{Title: "h0i!"})
}

func main() {
	rootdir, err := os.Getwd()
	if err != nil {
	rootdir = "no directory found"
	}
	http.Handle("/public/", http.StripPrefix("/public",
		http.FileServer(http.Dir(path.Join(rootdir, "public/")))))
	http.HandleFunc("/", mainHandler)
	
	//Listen on port 8080
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}