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
									"./app/views/profile.html", "./app/views/navbar.html", "./app/views/rolepage.html", "./app/views/projects.html", "./app/views/editProfile.html", "./app/views/projectPage.html", "./app/views/addWork.html", "./app/views/contestMain.html", "./app/views/submitCasting.html"))

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

func editProfileHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "editProfile", &Page{Title: "Edit Profile"})
}
func projectPageHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "projectPage", &Page{Title: "Project Page"})
}
func addWorkHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "addWork", &Page{Title: "Add Work"})
}
func contestMainHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "contestMain", &Page{Title: "Contest"})
}

func submitCastingHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "submitCasting", &Page{Title: "Submit Casting"})
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
	http.HandleFunc("/editProfile/", editProfileHandler)
	http.HandleFunc("/projectPage/", projectPageHandler)
	http.HandleFunc("/addWork/", addWorkHandler)
	http.HandleFunc("/contestMain/", contestMainHandler)
	http.HandleFunc("/submitCasting/", submitCastingHandler)
	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}