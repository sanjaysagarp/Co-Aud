package main

import (
	"html/template"
	"net/http"
	"fmt"
	"os"
	"path"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles("./app/views/header.html", "./app/views/footer.html", "./app/views/main.html", 
									"./app/views/profile.html", "./app/views/navbar.html", "./app/views/rolepage.html", "./app/views/projects.html"))

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

func profileHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "profile", &Page{Title: "Profile"})
}

func rolePageHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "rolepage", &Page{Title: "Role"})
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "projects", &Page{Title: "Projects"})
}

func main() {
	rootdir, err := os.Getwd()
	if err != nil {
		rootdir = "no directory found"
	}
	http.Handle("/public/", http.StripPrefix("/public",
		http.FileServer(http.Dir(path.Join(rootdir, "public/")))))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.HandleFunc("/role/", rolePageHandler)
	http.HandleFunc("/projects/", projectsHandler)

	//Listen on port 80
	fmt.Println("Server is listening on port 80...")
	http.ListenAndServe(":80", nil)
}