package main

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"encoding/json"
	"github.com/sanjaysagarp/Co-Aud/packages/user"
	"github.com/aaudis/GoRedisSession"
	"fmt"
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

var templates = template.Must(template.ParseFiles("./app/views/header.html", "./app/views/footer.html", "./app/views/main.html", 
									"./app/views/profile.html", "./app/views/navbar.html", "./app/views/rolepage.html", "./app/views/projects.html", "./app/views/editProfile.html", "./app/views/projectPage.html", "./app/views/addWork.html", "./app/views/contestMain.html", "./app/views/submitCasting.html", "./app/views/information.html", "./app/views/castings.html", "./app/views/home.html"))
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
	
	redis_session *rsess.SessionConnect
)

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

//The handlers.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "home", &Page{Title: "Home!"})
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
func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func castingsHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "castings", &Page{Title: "Casting List"})
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
	
	currentUser = user.FindUser(gUser.Email)
	if(currentUser == nil) {
		newUser := user.NewUser(gUser.Email, gUser.Name)
		user.InsertUser(newUser)
		currentUser = newUser
	}
	
	//get session and set whatever variables you want to access!
	s := redis_session.Session(w, r) // use this to retrieve current user session
	s.Set("DisplayName", currentUser.Email)
	s.Set("Email", currentUser.Email)
	s.Set("ID", currentUser.ID.String())
	fmt.Println(s.Get("Email"))
	fmt.Println(s.Get("ID"))
	//fmt.Fprintf(w, "Setting session variable done!")
	
	
	http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
	
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "information", &Page{Title: "Profile", Data: currentUser})

}

func main() {
	
	rsess.Prefix = "sess:" // session prefix (in Redis)
	rsess.Expire = 3600    // 60 minute session expiration

	// Connecting to Redis and creating storage instance
	temp_sess, err := rsess.New("sid", 0, "127.0.0.1", 27018)
	if err != nil {
		log.Printf("%s", err)
	}
	
	redis_session = temp_sess
	
	
	
	rootdir, err := os.Getwd()
	if err != nil {
		rootdir = "no directory found"
	}
	http.Handle("/public/", http.StripPrefix("/public",
		http.FileServer(http.Dir(path.Join(rootdir, "public/")))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.HandleFunc("/role/", rolePageHandler)
	http.HandleFunc("/projects/", projectsHandler)
	http.HandleFunc("/editProfile/", editProfileHandler)
	http.HandleFunc("/projectPage/", projectPageHandler)
	http.HandleFunc("/addWork/", addWorkHandler)
	http.HandleFunc("/contestMain/", contestMainHandler)
	http.HandleFunc("/submitCasting/", submitCastingHandler)
	http.HandleFunc("/login", googleLoginHandler)
	http.HandleFunc("/GoogleCallback", googleCallbackHandler)
	http.HandleFunc("/castings/", castingsHandler)
	http.HandleFunc("/user/", userHandler)
	

	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}