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
	// "reflect"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws/session"
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

}

func uploadTestHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "information", &Page{Title: "Upload pls!", Data: data})
}
//==========================================================================================================================================================

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
	data := setDefaultData(w, r)
	submission := make(map[string]interface{})
	// traits := strings.Split(r.FormValue("traits"), " ")
	
	submission["castEmail"] = r.FormValue("castEmail")
	submission["title"] = r.FormValue("title")
	submission["URL"] = r.FormValue("URL")
	submission["shortDescription"] = r.FormValue("shortDescription")
	submission["description"] = r.FormValue("description")
	// submission["cast"] = r.FormValue("cast")
	
	// submission["script"] = r.FormValue("script")
	// submission["deadline"] = r.FormValue("deadline")	
	// submission["gender"] = r.FormValue("gender")
	// submission["age"] = r.FormValue("age")
	// submission["traits"] = traits
	
	 newWork := work.NewWork(r.FormValue("title"), r.FormValue("URL"),r.FormValue("shortDescription"), r.FormValue("description"), r.FormValue("castEmail"), "seanyy")
	fmt.Println(newWork)
	fmt.Println(r.FormValue("castEmail[]"))
	
	work.InsertWork(newWork)

	data["form"] = r.Form
	data["submission"] = submission
	display(w, "seanTest", &Page{Title: "LULULULU", Data: data})
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

	w.Write([]byte("updated"))
	
}

//Submits an audition in auditions/{auditionid}
func submitAuditionHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	roleID := r.FormValue("id")
	err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        log.Printf("%s", err)
    }
	
	file, handler, err := r.FormFile("auditionFile")
	if(err != nil) {
		fmt.Printf("err opening file1: %s", err)
	}
	defer file.Close()
	
	
	fmt.Println("BLAH " + roleID )
	attachmentUrl := "/media/" + s.Get("Email") + "/" + handler.Filename
	
	uploader := s3manager.NewUploader(session.New())
    result, err := uploader.Upload(&s3manager.UploadInput{
        Body:   file,
        Bucket: aws.String("coaud"),
        Key:    aws.String(attachmentUrl),
    })
	
	if err != nil {
        log.Fatalln("Failed to upload", err)
    }

    log.Println("Successfully uploaded to", result.Location)
	
	//create a new audition and add the link
	audition := role.NewAudition(s.Get("Email"), attachmentUrl)
	curRole := role.FindRole(roleID)
	role.InsertAudition(audition, curRole)
	
	w.Write([]byte("uploaded"))
	
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
	
	a := role.FindRole(roleID)
	fmt.Println(a)
	
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
	http.HandleFunc("/upload/", uploadTestHandler)
	http.HandleFunc("/logout/", logoutHandler)
	
	//update handlers
	http.HandleFunc("/api/v1/updateUser/", updateUserHandler)
	http.HandleFunc("/api/v1/publishCasting/", publishCastingHandler)
	http.HandleFunc("/api/v1/submitComment/", submitCommentHandler)
	http.HandleFunc("/api/v1/submitAudition/", submitAuditionHandler)


	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}