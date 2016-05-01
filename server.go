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
	"github.com/sanjaysagarp/Co-Aud/packages/role"
	"github.com/aaudis/GoRedisSession"
	"fmt"
	"strings"
	"time"
  "reflect"
  "strconv"
)

//A Page structure
type Page struct {
	Title string
	Data map[string]interface{}
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

var templates = template.Must(template.ParseGlob("./app/views/*.html"))
var configFile, _ = ioutil.ReadFile("./secret/config.json")

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
	data := setDefaultData(w, r)
	display(w, "home", &Page{Title: "Home!", Data: data})
}

//FOR TESTING PURPOSES ONLY=================================================================================================================================
func theoTestPageHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data["currentUser"] = currentUser
	submission := make(map[string]interface{})


	fmt.Println(r.Form)

// TESTING CASTING SUBMISSION ==============================================================================================================================
	traits := strings.Split(r.FormValue("traits"), " ")

  layout := "2006-01-02"

	deadline, err := time.Parse(layout, r.FormValue("deadline"))
	if err != nil {
      fmt.Println(err)
      return
  }

  age, err := strconv.Atoi(r.FormValue("age"))
  if err != nil {
      fmt.Println(err)
      return
  }
	newRole := role.NewRole(r.FormValue("title"), currentUser.Email, r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"))
	fmt.Println(newRole)
	role.InsertRole(newRole)
	// getting string values
	submission["title"] = r.FormValue("title")
	submission["description"] = r.FormValue("description")
	submission["script"] = r.FormValue("script")
	submission["deadline"] = deadline
	submission["gender"] = r.FormValue("gender")
	submission["age"] = r.FormValue("age")
	submission["traits"] = traits


	// get picture
	// r.ParseMultipartForm(32 << 20)
 //  file, handler, err := r.FormFile("photo")
 //  if err != nil {
 //      fmt.Println(err)
 //      return
 //  }
 //  defer file.Close()
 //  fmt.Fprintf(w, "%v", handler.Header)
 //  filepathname := "C:/Users/Theo/Pictures/theo_test"+handler.Filename
 //  f, err := os.OpenFile(filepathname, os.O_WRONLY|os.O_CREATE, 0666)
 //  if err != nil {
 //      fmt.Println(err)
 //      return
 //  }
 //  fmt.Println(filepathname)
 //  defer f.Close()
 //  io.Copy(f, file)

// TESTING PROJECT SUBMISSION ===========================================
	// newWork = work.NewWork(r.FormValue("title"), r.FormValue("url"), r.FormValue("shortDescription"), r.FormValue("description"), cast []Cast, data["currentUser"].Email)
	// work.InsertWork(newWork)

	// submission["castMember"] = r.FormValue("castMember")
	// submission["castType"] = r.FormValue("castType")
	
	data["form"] = r.Form
	data["submission"] = submission

	// fmt.Println(submission)

	display(w, "theoTestPage", &Page{Title: "Theo Test", Data: data})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "profile", &Page{Title: "Profile", Data: data})
}

func rolePageHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	roleId := r.URL.Query().Get("id")
	fmt.Println(roleId)
	fmt.Println(reflect.TypeOf(roleId))
	role := role.FindRole(roleId)
	data["role"] = role
	data["postedBy"] = user.FindUser(role.UserEmail).DisplayName
	display(w, "rolepage", &Page{Title: "Role", Data: data})
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "projects", &Page{Title: "Projects", Data: data})
}

func editProfileHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "editProfile", &Page{Title: "Edit Profile", Data: data})
}

func projectPageHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "projectPage", &Page{Title: "Project Page", Data: data})
}

func addWorkHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "addWork", &Page{Title: "Add Work", Data: data})
}

func contestMainHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "contestMain", &Page{Title: "Contest", Data: data})
}

func submitCastingHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "submitCasting", &Page{Title: "Submit Casting", Data: data})
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func castingsHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	roleList := role.FindRoles()
	fmt.Println(roleList)
	fmt.Println(reflect.TypeOf(roleList))
	data["roles"] = roleList
	display(w, "castings", &Page{Title: "Casting List", Data: data})
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
	s.Set("ID", currentUser.Id.String())
	
	
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	
}

func setDefaultData(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data["currentUser"] = currentUser
	return data
}

func main() {
	
	rsess.Prefix = "sess:" // session prefix (in Redis)
	rsess.Expire = 3600    // 60 minute session expiration

	// Connecting to Redis and creating storage instance
	temp_sess, err := rsess.New("sid", 0, "127.0.0.1", 6379)
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
	http.HandleFunc("/theoTestPage/", theoTestPageHandler)

	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}