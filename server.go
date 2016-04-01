package main

import (
	"html/template"
	"net/http"
	"fmt"
	"os"
	"path"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
)

//A Page structure
type Page struct {
	Title string
}

//config structure
// type Config struct {
// 	Domain string
// 	Port int
// 	DB string
// 	SessionSecret string
// 	CookieSecret string
// }

//Compile templates on start
var templates = template.Must(template.ParseFiles("./app/views/header.html", "./app/views/footer.html", "./app/views/main.html"))
//var configFile, _ = ioutil.ReadFile("./secret/config.json")



var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:    "http://localhost:8080/GoogleCallback",
		ClientID:     os.Getenv("GOOGLEKEY"),
		ClientSecret: os.Getenv("GOOGLESECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile",
					"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
// Some random string, random for each request
	oauthStateString = "random"
)

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

//The handlers.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "main", &Page{Title: "h0i!"})
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}


func main() {
	rootdir, err := os.Getwd()
	if err != nil {
		rootdir = "no directory found"
	}
	http.Handle("/public/", http.StripPrefix("/public",
		http.FileServer(http.Dir(path.Join(rootdir, "public/")))))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/login", googleLoginHandler)
	http.HandleFunc("/GoogleCallback", googleCallbackHandler)
	
	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}