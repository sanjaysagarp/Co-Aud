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
	"github.com/sanjaysagarp/Co-Aud/packages/work"
	"github.com/aaudis/GoRedisSession"
	"fmt"
	"strings"
	"time"
	"math"
	//"reflect"
	"strconv"
	"gopkg.in/mgo.v2/bson"
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
		ClientID:     "688463917821-p5u7nvg7eovcjr92o7e8986b3tl3qdlr.apps.googleusercontent.com",
		ClientSecret: "nyIHJVB8Fzx76kSL7SMFFRFP",
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

}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	userID := r.URL.Query().Get("id")
	user := user.FindUserById(userID)
	postedRoles, rolesCount := role.FindRoles(bson.M{"useremail": user.Email}, 0, 3)
	data["user"] = user
	data["postedRoles"] = postedRoles
	data["rolesCount"] = rolesCount
	display(w, "profile", &Page{Title: user.DisplayName, Data: data})
}

func rolePageHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	roleID := r.URL.Query().Get("id")
	role := role.FindRole(roleID)
	data["role"] = role
	data["author"] = user.FindUser(role.UserEmail)
	fmt.Println(role.Comment)
	display(w, "rolepage", &Page{Title: "Role", Data: data})
}
func workHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data := setDefaultData(w, r)
	workId := r.URL.Query().Get("id")
	work := work.FindWork(workId)
	URL := work.URL
	youtubeCode := strings.Split(URL, "=")
	//fmt.Println("check below this line")
	fmt.Println(youtubeCode)
	data["youtubeCode"] = youtubeCode[1] // causes error,need to find better way
	data["work"] = work
	data["user"] = currentUser
	//data["author"] = user.FindUser(work.UserEmail)
	//fmt.Println(role.Comment)
	display(w, "seanTest", &Page{Title: "Work", Data: data})
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	projectAmount := 16
	
	//pagination
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page"))
	data["currentPage"] = pageNumber
	if err != nil {
		fmt.Println(err)
	}
	//zero index page number for skip calculation when querying mongo
	if pageNumber != 0 {
		pageNumber --
	}
	//get roles
	projectList, projectCount := work.FindWorks(nil, (pageNumber)*projectAmount, projectAmount)
	
	//get max page number
	maxPage := int(math.Ceil(float64(projectCount)/float64(projectAmount)))
	
	data["works"] = projectList
	data["workCount"] = projectCount
	data["workAmount"] = projectAmount
	data["maxPage"] = maxPage
	
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

func publishCastingHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))

	fmt.Println(r.Form)

	// converting data to valid format
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
	  
	// adding new role into db
	roleID := bson.NewObjectId()
	newRole := role.NewRole(r.FormValue("title"), currentUser.Email, r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"), roleID)
	role.InsertRole(newRole)

	urlParts := []string{"/role/?id=", roleID.Hex()}
	url := strings.Join(urlParts, "")
	// redirect to role page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func castingsHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	roleAmount := 16
	
	//pagination
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page"))
	data["currentPage"] = pageNumber
	if err != nil {
		fmt.Println(err)
	}
	//zero index page number for skip calculation when querying mongo
	if pageNumber != 0 {
		pageNumber --
	}
	//get roles
	roleList, rolesCount := role.FindRoles(nil, (pageNumber)*roleAmount, roleAmount)
	
	//get max page number
	maxPage := int(math.Ceil(float64(rolesCount)/float64(roleAmount)))
	
	data["roles"] = roleList
	data["rolesCount"] = rolesCount
	data["roleAmount"] = roleAmount
	data["maxPage"] = maxPage
	
	display(w, "castings", &Page{Title: "Casting List", Data: data})
}

func seanTestHands(w http.ResponseWriter, r *http.Request) {
	//data := setDefaultData(w, r)
	//submission := make(map[string]interface{})

	castsAttendees := strings.Split(r.FormValue("castList"), ",")
	castRoles := strings.Split(r.FormValue("castRoles"), ",")
	fmt.Println(len(castsAttendees))
	fmt.Println(len(castRoles))

	castContainer := make([]work.Cast, 0)
	for i := 0; i < len(castsAttendees); i++ {
		fmt.Println(castsAttendees[i])
		fmt.Println(castRoles[i])
		castUser := user.FindUser(castsAttendees[i])
		newCast := work.NewCast(castUser, castRoles[i])
		castContainer = append(castContainer, newCast)
	}
	
	// data["castEmail"] = r.FormValue("castEmail")
	// data["title"] = r.FormValue("title")
	// data["URL"] = r.FormValue("URL")
	// data["shortDescription"] = r.FormValue("shortDescription")
	// data["description"] = r.FormValue("description")
	// data["castHolder"] = castContainer
	projectId := bson.NewObjectId()
	//s := redis_session.Session(w, r)
	
	newWork := work.NewWork(r.FormValue("title"), r.FormValue("url"),r.FormValue("shortDescription"), r.FormValue("description"), castContainer, currentUser, projectId)
	work.InsertWork(newWork)
	fmt.Println(newWork)
	// data["form"] = r.Form
	// data["submission"] = submission
	
	//w.Write([]byte("updated"))
	//display(w, "seanTest", &Page{Title: "LULULULU", Data: data})
	urlParts := []string{"/work/?id=", projectId.Hex()}
	url := strings.Join(urlParts, "")
	// redirect to role page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}


// INFINITE SCROLL STUFF GOES HERE; NOT COMPLETE
// func getMoreCastingsHandler()

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
	s.Set("ID", currentUser.Id.Hex()) //gives hex value of id for access
	
	
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	s.Destroy(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	//need to get form fields web page
	err := r.ParseForm()
    if err != nil {
        log.Printf("%s", err)      
    }
	
	//create a new user struct
	editedUser := user.NewChangeUser(r.FormValue("displayName"), r.FormValue("title"), r.FormValue("aboutMe"), r.FormValue("personalWebsite"),  r.FormValue("facebookURL"), r.FormValue("twitterURL"), r.FormValue("instagramURL"))
	
	user.UpdateUser(s.Get("ID"), editedUser)
	
	//need to write to page for ajax call
	w.Write([]byte("updated"))
	
}

func submitCommentHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	
	collection := r.FormValue("collection")
	message := r.FormValue("content")
	roleID := r.FormValue("id")
	
	
	newComment := role.NewComment(currentUser, message)
	curRole := role.FindRole(roleID)

	role.InsertComment(curRole.Comment, newComment, collection, curRole.Id.Hex())
	
	w.Write([]byte("updated"))

}

func setDefaultData(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	data := make(map[string]interface{})
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	
	data["currentUser"] = currentUser
	if(currentUser != nil) {
		s.Set("DisplayName", currentUser.Email)
		s.Set("Email", currentUser.Email)
		s.Set("ID", currentUser.Id.Hex()) 
	}
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
	http.HandleFunc("/work/", workHandler)
	http.HandleFunc("/projects/", projectsHandler)
	http.HandleFunc("/profile/edit/", editProfileHandler)
	http.HandleFunc("/projectPage/", projectPageHandler)
	http.HandleFunc("/addWork/", addWorkHandler)
	http.HandleFunc("/contestMain/", contestMainHandler)
	http.HandleFunc("/submitCasting/", submitCastingHandler)
	http.HandleFunc("/login", googleLoginHandler)
	http.HandleFunc("/GoogleCallback", googleCallbackHandler)
	http.HandleFunc("/castings/", castingsHandler)
	http.HandleFunc("/theoTestPage/", theoTestPageHandler)
	http.HandleFunc("/logout/", logoutHandler)
	//http.HandleFunc("/seanTest/", seanTestHands)
	
	//update handlers
	http.HandleFunc("/api/v1/updateUser/", updateUserHandler)
	http.HandleFunc("/api/v1/publishCasting/", publishCastingHandler)
	http.HandleFunc("/api/v1/submitComment/", submitCommentHandler)
	http.HandleFunc("/api/v1/publishWork/", seanTestHands)
	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}