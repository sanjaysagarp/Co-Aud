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
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/service/s3"
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
	
	//get roles
	roleList, totalRolesCount := role.FindRoles(nil, 0, 16)
	projectList, totalProjectList := work.FindWorks(nil, 0, 6)
	
	fmt.Println("total number of projects: ", totalProjectList)
	data["totalRolesCount"] = totalRolesCount
	data["roles"] = roleList
	data["projects"] = projectList
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
	postedRoles, rolesCount := role.FindRoles(bson.M{"user.email": user.Email}, 0, 3)
	postedWorks, workCount := work.FindWorks(bson.M{"user.email": user.Email}, 0, 3)
	data["user"] = user
	data["postedRoles"] = postedRoles
	data["rolesCount"] = rolesCount
	data["postedWorks"] = postedWorks
	data["workCount"] = workCount
	display(w, "profile", &Page{Title: user.DisplayName, Data: data})
}

func rolePageHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	roleID := r.URL.Query().Get("id")
	role := role.FindRole(roleID)
	data["role"] = role
	
	// svc := s3.New(session.New())
	// //Pre-signs all audio clips so they cannot be downloaded! -- Do we want this?
	// for _,audition := range role.Audition {
	// 	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
	// 		Bucket: aws.String("coaud"),
	// 		Key:    aws.String(audition.AttachmentUrl),
	// 	})
	// 	temp, err := req.Presign(15 * time.Minute)
	// 	if err != nil {
	// 		log.Println("Failed to sign request", err)
	// 	}
	// 	audition.TempUrl = temp
	// }
	
	data["author"] = user.FindUser(role.User.Email)
	fmt.Println(role.Comment)
	display(w, "rolepage", &Page{Title: role.Title, Data: data})
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data := setDefaultData(w, r)
	workID := r.URL.Query().Get("id")
	work := work.FindWork(workID)
	//fmt.Println("check below this line")
	data["youtubeCode"] = work.GetYoutubeID() // causes error,need to find better way
	data["work"] = work
	data["user"] = currentUser
	//data["author"] = user.FindUser(work.UserEmail)
	//fmt.Println(role.Comment)
	display(w, "seanTest", &Page{Title: "Work", Data: data})
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	projectAmount := 16
	pageAmount := 5
	
	//pagination
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page")) //used for getting roles
	currentPage := pageNumber //used for getting page list
	if currentPage <= 0 {
		currentPage = 1
	}
	data["currentPage"] = currentPage
	data["prevPage"] = currentPage - 1
	data["nextPage"] = currentPage + 1
	if err != nil {
		fmt.Println(err)
	}
	//zero index page number for skip calculation when querying mongo
	if pageNumber != 0 {
		pageNumber --
	}
	
	//get projects
	projectList, projectCount := work.FindWorks(nil, (pageNumber)*projectAmount, projectAmount)
	fmt.Println(projectList)
	//more params for pagination
	maxPage := int(math.Ceil(float64(projectCount)/float64(projectAmount)))
	pageList := getPageList(maxPage, currentPage, pageAmount)
	
	data["works"] = projectList
	data["pageList"] = pageList
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
	newRole := role.NewRole(r.FormValue("title"), currentUser, r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"), roleID)
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
	//params for pagination
	roleAmount := 16
	pageAmount:= 5
	
	//pagination
	pageNumber, err := strconv.Atoi(r.URL.Query().Get("page")) //used for getting roles
	currentPage := pageNumber //used for getting page list
	if currentPage <= 0 {
		currentPage = 1
	}
	data["currentPage"] = currentPage
	data["prevPage"] = currentPage - 1
	data["nextPage"] = currentPage + 1
	if err != nil {
		fmt.Println(err)
	}
	//zero index page number for skip calculation when querying mongo
	if pageNumber != 0 {
		pageNumber --
	}
	
	//get roles
	roleList, rolesCount := role.FindRoles(nil, (pageNumber)*roleAmount, roleAmount)
	
	//more params for pagination
	maxPage := int(math.Ceil(float64(rolesCount)/float64(roleAmount)))
	pageList := getPageList(maxPage, currentPage, pageAmount)
	
	data["roles"] = roleList
	data["maxPage"] = maxPage
	data["pageList"] = pageList
	display(w, "castings", &Page{Title: "Casting List", Data: data})
}


func seanTestHands(w http.ResponseWriter, r *http.Request) {
	//data := setDefaultData(w, r)
	//submission := make(map[string]interface{})
	fmt.Println(r.FormValue("castList"))
	fmt.Println(r.FormValue("castRoles"))
	castsAttendees := strings.Split(r.FormValue("castList"), ",")
	castRoles := strings.Split(r.FormValue("castRoles"), ",")
	fmt.Println(len(castsAttendees))
	fmt.Println(len(castRoles))

	castContainer := make([]work.Cast, 0)
	for i := 0; i < len(castsAttendees); i++ {
		castUser := user.FindUser(castsAttendees[i])
		newCast := work.NewCast(castUser, castRoles[i])
		castContainer = append(castContainer, newCast)
	}
	
	projectId := bson.NewObjectId()
	//s := redis_session.Session(w, r)
	newWork := work.NewWork(r.FormValue("title"), r.FormValue("URL"),r.FormValue("shortDescription"), r.FormValue("description"), castContainer, currentUser, projectId)
	work.InsertWork(newWork)
	fmt.Println(newWork)
	//display(w, "seanTest", &Page{Title: "LULULULU", Data: data})
	urlParts := []string{"/work/?id=", projectId.Hex()}
	url := strings.Join(urlParts, "")
	// redirect to project page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
//gets the numbers of the pages that will be shown in pagination given the max page, current page,
//and the amount of pages you want displayed
func getPageList(maxPage int, curPage int, amount int) []int{
	var result []int
	var min int
	var max int
	
	if (curPage - (amount/2) <= 1) { //first few pages
		min = 1
		if (maxPage > (curPage + amount - 1)) { //if there are more pages than what we will show
			//get as many pages 1 to amount
			max = amount
		} else { //the amount of pages total is less than or equal to the max number of pages
			//get the pages from 1 to max page
			max = maxPage
		}
	} else if (curPage + (amount/2) >= maxPage) { //last few pages
		//get as many pages maxPage - (amount - 1) to maxPage
		min = maxPage - (amount - 1)
		max = maxPage
	} else { //somewhere in the middle
		//get current page - amount/2 to current page - amount/2 + (amount-1)
		min = curPage - (amount/2)
		max = min + (amount - 1)
	}
	
	for i := min; i <= max; i++ {
		result = append(result, i)
	}
	fmt.Println(result)
	return result

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
	if err != nil {
		return
	}

	bytes, err := file.Seek(0,2)
	if(err != nil) {
		panic(err)
	}
	
	var kilobytes int64
	kilobytes = (bytes / 1024)
	
	var megabytes float64
	megabytes = (float64)(kilobytes / 1024)
	
	if(megabytes < 6) {
		attachmentURL := "/roles/" + roleID + "/" + s.Get("Email") + "/" + handler.Filename
	
		uploader := s3manager.NewUploader(session.New())
		result, err := uploader.Upload(&s3manager.UploadInput{
			Body:   file,
			Bucket: aws.String("coaud"),
			Key:    aws.String(attachmentURL),
		})
		
		if err != nil {
			log.Fatalln("Failed to upload", err)
		}

		log.Println("Successfully uploaded to", result.Location)
		
		//create a new audition and add the link
		audition := role.NewAudition(user.FindUser(s.Get("Email")), result.Location)
		curRole := role.FindRole(roleID)
		role.InsertAudition(audition, curRole)
		
		w.Write([]byte("uploaded"))
	} else {
		w.Write([]byte("rejected"))
	}
	
	
	
}


func submitRoleCommentHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	
	collection := "roles"
	message := r.FormValue("content")
	roleID := r.FormValue("id")
	
	
	newComment := role.NewComment(currentUser, message)
	curRole := role.FindRole(roleID)

	role.InsertComment(curRole.Comment, newComment, collection, curRole.Id.Hex())
	
	a := role.FindRole(roleID)
	fmt.Println(a)
	
	w.Write([]byte("updated"))
}

func getRoleHandler(w http.ResponseWriter, r *http.Request) {
	//need to get form fields web page
	page, err := strconv.Atoi(r.FormValue("page"))
	fmt.Println(page)
	if(err != nil) {
		panic(err)
	}
	
	//create a new user struct
	roles, rolesCount := role.FindRoles(nil, page*16, 16)
	fmt.Println(rolesCount)
	var html string
	if len(roles) != 0 {
		for index, role := range roles {
			html += "<div class=\"thumbnail col-sm-3 col-lg-3 col-xs-6 col-md-3\"><a href=\"/role/?id=" + role.Id.Hex() + "\"><h2>" + role.Title + "</h2><img class=\"img-responsive\" src=\"/public/img/placeholder.png\"></a></div>"
			fmt.Println(index)
		}
	}
	w.Write([]byte(html))
	
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
	http.HandleFunc("/upload/", uploadTestHandler)
	http.HandleFunc("/logout/", logoutHandler)
	//http.HandleFunc("/seanTest/", seanTestHands)
	
	//update handlers
	http.HandleFunc("/api/v1/updateUser/", updateUserHandler)
	http.HandleFunc("/api/v1/publishCasting/", publishCastingHandler)
	http.HandleFunc("/api/v1/submitAudition/", submitAuditionHandler)
	http.HandleFunc("/api/v1/submitComment/", submitRoleCommentHandler)
	http.HandleFunc("/api/v1/publishWork/", seanTestHands)
	http.HandleFunc("/api/v1/submitRoleComment/", submitRoleCommentHandler)
	http.HandleFunc("/api/v1/getRole/", getRoleHandler)

	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}