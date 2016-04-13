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
	"encoding/json"
	"github.com/sanjaysagarp/Co-Aud/packages"
)

//A Page structure
type Page struct {
	Title string
	Data interface{}
}

//UserPage struct
type UserPage struct {
	Title string
}

//GoogleUser struct that captures initial user information for acct creation
type GoogleUser struct {
	Email string `json:"email"`
	Name string `json:"name"`
}

//Compile templates on start
var templates = template.Must(template.ParseFiles("./app/views/header.html", "./app/views/footer.html", "./app/views/main.html", "./app/views/information.html"))
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
	currentUser = &user.User{}
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
	var gUser GoogleUser
	json.Unmarshal(contents, &gUser)
	//fmt.Println(gUser.Email)
	
	//need to search for email in our db -> if found, navigate back to homepage?
	
	currentUser = user.FindUser(gUser.Email)
	if(currentUser == nil) {
		newUser := user.NewUser(gUser.Email, gUser.Name)
		user.InsertUser(newUser)
		currentUser = newUser
	}
	
	
	http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
	
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	//should check if token is still valid?
	display(w, "information", &Page{Title: "Profile", Data: currentUser})
}
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	//Checks if user has a valid oauth token
	// code := r.FormValue("code")
	// token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	// if err != nil {
	// 	fmt.Println("Code exchange failed with '%s'\n", err)
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	// if(!token.Valid()) {
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// }
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
	//http.HandleFunc("/createUser", createUserHandler)
	http.HandleFunc("/user", userHandler)
	
	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}