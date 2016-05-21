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
	"github.com/sanjaysagarp/Co-Aud/packages/project"
	"github.com/aaudis/GoRedisSession"
	"fmt"
	"strings"
	"time"
	"math"
	//"reflect"
	"strconv"
	"gopkg.in/mgo.v2"
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
	
	//get roles
	roleList, totalRolesCount := role.FindRoles(nil, 0, 16)
	projectList, totalProjectList := project.FindProjects(nil, 0, 6)
	
	fmt.Println("total number of projects: ", totalProjectList)
	data["totalRolesCount"] = totalRolesCount
	data["roles"] = roleList
	data["projects"] = projectList
	display(w, "home", &Page{Title: "Home!", Data: data})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	userID := r.URL.Query().Get("id")
	user := user.FindUserById(userID)
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	postedRoles, rolesCount := role.FindRoles(bson.M{"user": dbRefUser}, 0, 3)
	postedProjects, projectCount := project.FindProjects(bson.M{"user": dbRefUser}, 0, 3)
	data["user"] = user
	data["postedRoles"] = postedRoles
	data["rolesCount"] = rolesCount
	data["postedProjects"] = postedProjects
	data["projectCount"] = projectCount
	display(w, "viewProfile", &Page{Title: user.DisplayName, Data: data})
}
func editAuditionHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	userID := r.URL.Query().Get("id")
	user := user.FindUserById(userID)
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	postedRoles, rolesCount := role.FindRoles(bson.M{"user": dbRefUser}, 0, 500)
	data["user"] = user
	data["postedRoles"] = postedRoles
	data["rolesCount"] = rolesCount
	display(w, "editAudition", &Page{Title: user.DisplayName, Data: data})
}

func editProjectHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	userID := r.URL.Query().Get("id")
	user := user.FindUserById(userID)
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	postedProjects, projectCount := project.FindProjects(bson.M{"user": dbRefUser}, 0, 500)
	data["user"] = user
	data["postedProjects"] = postedProjects
	data["projectCount"] = projectCount
	display(w, "editProject", &Page{Title: user.DisplayName, Data: data})
}

func changeProjectHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	projectId := r.URL.Query().Get("id")
	project := project.FindProject(projectId)
	data["project"] = project
	data["castLength"] = len(project.Cast)
	display(w, "changeProject", &Page{Title: project.Title, Data: data})
}

func changeAuditionHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	roleId := r.URL.Query().Get("id")
	role := role.FindRole(roleId)
	data["role"] = role
	display(w, "changeAudition", &Page{Title: role.Title, Data: data})
}

func roleHandler(w http.ResponseWriter, r *http.Request) {
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
	
	data["author"] = role.GetUser()
	display(w, "viewRole", &Page{Title: role.Title, Data: data})
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data := setDefaultData(w, r)
	projectID := r.URL.Query().Get("id")
	project := project.FindProject(projectID)
	//fmt.Println("check below this line")
	data["youtubeCode"] = project.GetYoutubeID() // causes error,need to find better way
	data["project"] = project
	data["user"] = currentUser
	//data["author"] = user.FindUser(project.UserEmail)
	//fmt.Println(role.Comment)
	display(w, "viewProject", &Page{Title: "Project", Data: data})
}

func teamHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data := setDefaultData(w, r)
	teamID := r.URL.Query().Get("id")
	//fmt.Println(teamID)
	team := role.FindTeam(teamID)
	fmt.Println(team)
	data["team"] = team
	data["user"] = currentUser
	display(w, "viewTeam", &Page{Title: "Team", Data: data})
}


func projectBrowseHandler(w http.ResponseWriter, r *http.Request) {
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
	projectList, projectCount := project.FindProjects(nil, (pageNumber)*projectAmount, projectAmount)
	fmt.Println(projectList)
	//more params for pagination
	maxPage := int(math.Ceil(float64(projectCount)/float64(projectAmount)))
	pageList := getPageList(maxPage, currentPage, pageAmount)
	
	data["projects"] = projectList
	data["pageList"] = pageList
	data["maxPage"] = maxPage
	
	display(w, "browseProjects", &Page{Title: "Projects", Data: data})
}

func browseTeamsHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	
	newContest := role.FindContests()
	contestList := role.FindContest(newContest[0].Id.Hex())
	
	fmt.Println(contestList)
	
	data["contests"] = contestList
	
	display(w, "browseTeams", &Page{Title: "Teams", Data: data})
}



func editProfileHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "editProfile", &Page{Title: "Edit Profile", Data: data})
}

// This was for testing purposes
// func projectPageHandler(w http.ResponseWriter, r *http.Request) {
// 	data := setDefaultData(w, r)
// 	display(w, "projectPage", &Page{Title: "Project Page", Data: data})
// }

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "createProject", &Page{Title: "Create Project", Data: data})
}

func createTeamHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "createTeam", &Page{Title: "Create Team", Data: data})
}


func contestMainHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	data := setDefaultData(w, r)
	
	// newContest := role.FindContests()
	// currentContest := role.FindContest(newContest[0].Id.Hex())
	
	contestId := r.URL.Query().Get("id")
	if (contestId == "") {
		newContest := role.FindContests()
		contestId = newContest[0].Id.Hex()
	} 
	contest := role.FindContest(contestId)
	data["contest"] = contest
	data["user"] = currentUser
	display(w, "viewContest", &Page{Title: "Contest", Data: data})
}


func createRoleHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "createRole", &Page{Title: "Create Role", Data: data})
}

func createContestHandler(w http.ResponseWriter, r *http.Request) {
	data := setDefaultData(w, r)
	display(w, "createContest", &Page{Title: "Create Contest", Data: data})
}

func updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	roleId := r.URL.Query().Get("id")
	rawTraits := strings.Split(r.FormValue("traits"), ",")
	var traits []string
	for i := 0; i < len(rawTraits); i++ {
		s := strings.TrimSpace(rawTraits[i])
		if (s != "") {
			traits = append(traits, s)
		}
	}
	layout := "2006-01-02"
	fmt.Println(r.FormValue("deadline"))
	deadline, err := time.Parse(layout, r.FormValue("deadline"))
	if err != nil {
		fmt.Println(err)
		return
 	 }
	fmt.Println("Have we gotten this far")
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		fmt.Println(err)
		return
	}
	file, handler, err := r.FormFile("photo")
	if (err == nil) {
		fmt.Println("success photo is in there")
		defer file.Close()
		bytes, err := file.Seek(0,2)
		if err != nil {
			panic(err)
		}
		var kilobytes int64
		kilobytes = (bytes / 1024)
		
		var megabytes float64
		megabytes = (float64)(kilobytes / 1024)
		
		roleID := bson.NewObjectId()

		if(megabytes < 4) {
			attachmentURL := "/roles/" + roleID.Hex() + "/" + s.Get("Email") + "/" + handler.Filename
		
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
			
			editedRoleWithPhoto := role.ChangedRoleWithPhoto(r.FormValue("title"), r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"), result.Location)
			
			role.UpdateRoleWithPhoto(roleId, editedRoleWithPhoto)

		urlParts := []string{"/auditions/edit/?id=", currentUser.Id.Hex(), "#success"}
		url := strings.Join(urlParts, "")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			urlParts := []string{"/auditions/edit/?id=", currentUser.Id.Hex(), "#failure"}
			url := strings.Join(urlParts, "")
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		}
	} else {
		fmt.Println("No photo is in there")
		editedRoleWithNoPhoto := role.ChangedRoleNoPhoto(r.FormValue("title"), r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"))
		
		role.UpdateRoleNoPhoto(roleId, editedRoleWithNoPhoto)
		
		urlParts := []string{"/auditions/edit/?id=", currentUser.Id.Hex(), "#success"}
		url := strings.Join(urlParts, "")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
	
}


func submitRoleHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))

	fmt.Println(r.FormValue("photo"))
	// converting data to valid format
	rawTraits := strings.Split(r.FormValue("traits"), ",")
	var traits []string
	//clean traits
	for i := 0; i < len(rawTraits); i++ {
		s := strings.TrimSpace(rawTraits[i])
		if (s != "") {
			traits = append(traits, s)
		}
	}

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
	
	file, handler, err := r.FormFile("photo")
	
	if (err == nil) {
		defer file.Close()
		bytes, err := file.Seek(0,2)
		if err != nil {
			panic(err)
		}
		
		//get file size in kilobytes and megabytes
		var kilobytes int64
		kilobytes = (bytes / 1024)
		
		var megabytes float64
		megabytes = (float64)(kilobytes / 1024)
		
		// adding new role into db
		roleID := bson.NewObjectId()
		//TODO: add appropriate size limit
		// sss := strconv.FormatFloat(megabytes, 'f', 6, 64)
		// fmt.Println("Size: " + sss )
		if(megabytes < 4) {
			attachmentURL := "/roles/" + roleID.Hex() + "/" + s.Get("Email") + "/" + handler.Filename
		
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
			
			newRole := role.NewRole(r.FormValue("title"), currentUser, r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"), roleID, result.Location)
			role.InsertRole(newRole)

			// urlParts := []string{"/auditions/?id=", newRole.Id.Hex()}
			// url := strings.Join(urlParts, "")
			w.Write([]byte(newRole.Id.Hex()))
			//http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			//handle response if greater than 4 megabytes! -- NEED TO MAKE RESPONSIVE
			w.Write([]byte("rejected"))
		}
	} else {
		fmt.Printf("err opening image file: %s", err)
		fmt.Println("Placing default image..")
		roleID := bson.NewObjectId()

		newRole := role.NewRole(r.FormValue("title"), currentUser, r.FormValue("description"), r.FormValue("script"), deadline, traits, age, r.FormValue("gender"), roleID, "/public/img/default_role_pic.png")
		role.InsertRole(newRole)
		fmt.Println(newRole)
		// urlParts := []string{"/auditions/?id=", newRole.Id.Hex()}
		// url := strings.Join(urlParts, "")
		
		// http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		w.Write([]byte(newRole.Id.Hex()))
	}
	
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func roleBrowseHandler(w http.ResponseWriter, r *http.Request) {
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
	display(w, "browseAuditions", &Page{Title: "Roles", Data: data})
}

//api call for submitting project
func submitProjectHandler(w http.ResponseWriter, r *http.Request) {
	//data := setDefaultData(w, r)
	//submission := make(map[string]interface{})
	// **TOP
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	
	r.ParseForm()
	castsAttendees := r.Form["castEmail[]"]
	castRoles := r.Form["castRole[]"]

	var castContainer []*project.Cast
	var newCast *project.Cast
	for i := 0; i < len(castsAttendees); i++ {
		castId := bson.NewObjectId()
		castUser := user.FindUser(castsAttendees[i])
		//fmt.Println("ARWAREWARNEW " + castUser.Email)
		newCast = project.NewCast(castUser, castRoles[i], castId)
		project.InsertCast(newCast)
		castContainer = append(castContainer, newCast)
	}
	
	projectId := bson.NewObjectId()
	fmt.Println("ID :::: ", projectId)
	//s := redis_session.Session(w, r)
	newProject := project.NewProject(r.FormValue("title"), r.FormValue("url"),r.FormValue("shortDescription"), r.FormValue("description"), castContainer, currentUser, projectId)
	//fmt.Println("new project: ", newProject.Id)
	project.InsertProject(newProject)
	
	urlParts := []string{"/projects/?id=", projectId.Hex()}
	url := strings.Join(urlParts, "")
	// redirect to project page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	
	// ** Bottom
}

func submitTeamHandler(w http.ResponseWriter, r *http.Request) {
	//data := setDefaultData(w, r)
	//submission := make(map[string]interface{})
	// **TOP
	r.ParseForm()
	teamMembers := r.Form["teamEmails"]
	
	newContest := role.FindContests()
	currentContest := role.FindContest(newContest[0].Id.Hex())
	
	teamContainer := make([]*user.User, 0)
	for i := 0; i < len(teamMembers); i++ {
		newUser := user.FindUser(teamMembers[i])
		teamContainer = append(teamContainer, newUser)
	}
	
	fmt.Println("jk below this")
	fmt.Println(teamContainer)
	teamId := bson.NewObjectId()
	//s := redis_session.Session(w, r)
	
	// ****** CLEANING ********
	newTeam := role.NewTeam(teamContainer, r.FormValue("teamName"), r.FormValue("motto"), teamId)
	
	role.InsertNewTeam(newTeam)
	
	currentContest.InsertTeam(newTeam)
	fmt.Println("Below this")

	urlParts := []string{"/teams/?id=", teamId.Hex()}
	url := strings.Join(urlParts, "")
	// redirect to project page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func submitContestHandler(w http.ResponseWriter, r *http.Request) {

	contestId := bson.NewObjectId()
	//s := redis_session.Session(w, r)
	layout := "2006-01-02"

	deadline, err := time.Parse(layout, r.FormValue("deadline"))
	if err != nil {
		fmt.Println(err)
		return
	}
	  
	newContest := role.NewContest(r.FormValue("title"), r.FormValue("description"), r.FormValue("photo"), deadline, contestId)
	
	role.InsertContest(newContest)
	
	urlParts := []string{"/contest/?id=", contestId.Hex()}
	url := strings.Join(urlParts, "")
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

func updateProjectHandler(w http.ResponseWriter, r*http.Request) {
	//s := redis_session.Session(w, r)
	r.ParseForm()
	//currentUser := user.FindUser(s.Get("Email"))
	castsAttendees := r.Form["castEmail[]"]
	//fmt.Println(castsAttendees)
	castRoles := r.Form["castRole[]"]
	//fmt.Println(castRoles)
	projectId := r.URL.Query().Get("id")
	var castContainer []*project.Cast
	var newCast *project.Cast
	for i := 0; i < len(castsAttendees); i++ {
		castId := bson.NewObjectId()
		//fmt.Println(castId)
		castUser := user.FindUser(castsAttendees[i])
		//fmt.Println(castUser)
		newCast = project.NewCast(castUser, castRoles[i], castId)
		project.InsertCast(newCast)
		castContainer = append(castContainer, newCast)
	}
	editedProject := project.ChangedProject(r.FormValue("title"), r.FormValue("url"), r.FormValue("shortDescription"), r.FormValue("description"), castContainer)
	project.UpdateProject(projectId, editedProject)

	urlParts := []string{"/projects/?id=", projectId, "#success"}
	url := strings.Join(urlParts, "")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//Submits an audition in auditions/{auditionid}
func submitAuditionHandler(w http.ResponseWriter, r *http.Request) {
	s := redis_session.Session(w, r)
	fmt.Println(r.Form)
	err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        log.Printf("%s", err)
    }
	
	roleID := r.FormValue("id")
	fmt.Println(r.Form)
	file, handler, err := r.FormFile("auditionFile")
	defer file.Close()
	if err != nil {
		fmt.Printf("err opening audition file: %s", err)
		return
	}

	bytes, err := file.Seek(0,2)
	if err != nil {
		panic(err)
	}
	
	var kilobytes int64
	kilobytes = (bytes / 1024)
	
	var megabytes float64
	megabytes = (float64)(kilobytes / 1024)
	
	if(megabytes < 6) {
		attachmentURL := "/auditions/" + roleID + "/" + s.Get("Email") + "/" + handler.Filename
	
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
		auditionID := bson.NewObjectId()
		audition := role.NewAudition(user.FindUser(s.Get("Email")), result.Location, auditionID)
		curRole := role.FindRole(roleID)
		role.InsertAudition(audition, curRole)
		
		w.Write([]byte("uploaded"))
	} else {
		w.Write([]byte("rejected"))
	}
}



func submitRoleCommentHandler(w http.ResponseWriter, r *http.Request) {
	submitCommentHandler(w, r, "roles", true)
}

func submitAuditionCommentHandler(w http.ResponseWriter, r *http.Request) {
	submitCommentHandler(w, r, "auditions", false)
}

func submitCommentHandler(w http.ResponseWriter, r *http.Request, collection string, recentOrder bool) {
	s := redis_session.Session(w, r)
	currentUser := user.FindUser(s.Get("Email"))
	
	message := r.FormValue("content")
	id := r.FormValue("id")
	
	commentID := bson.NewObjectId()
	newComment := role.NewComment(currentUser, message, commentID)

	role.InsertComment(newComment, collection, id, recentOrder)
	
	w.Write([]byte(`<li class="comment-posted media"><div class="media-left"><a href="/profile/?id=` + currentUser.Id.Hex() + `"><img class="img-profile media-object img-circle" src="/public/img/default_profile_pic.png"></a></div><div class="media-body"><a href="/profile/?id=` + currentUser.Id.Hex() + `"><h4 class="media-heading">` + currentUser.DisplayName + `</h4></a><span class="timestamp">` + newComment.TimeStamp.Format("Jan 02, 2006") + `</span><p class="comment-message">` + message + `</p></div></li>`))
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
	//http.HandleFunc("/role/", rolePageHandler)
	http.HandleFunc("/teams/browse", browseTeamsHandler)
	http.HandleFunc("/projects/", projectHandler)
	http.HandleFunc("/projects/browse/", projectBrowseHandler)
	http.HandleFunc("/profile/edit/", editProfileHandler)
	//http.HandleFunc("/projectPage/", projectPageHandler)
	http.HandleFunc("/projects/create", createProjectHandler)
	http.HandleFunc("/contest/", contestMainHandler)
	http.HandleFunc("/teams/", teamHandler)
	http.HandleFunc("/teams/create", createTeamHandler)
	http.HandleFunc("/contest/create", createContestHandler)
	
	//Edit functions for auditions and projects //START
	http.HandleFunc("/project/edit/", editProjectHandler)
	http.HandleFunc("/auditions/edit/", editAuditionHandler)
	http.HandleFunc("/EditProject/", changeProjectHandler)
	http.HandleFunc("/EditAudition/", changeAuditionHandler)
	http.HandleFunc("/update/project/", updateProjectHandler)
	http.HandleFunc("/update/audition/", updateRoleHandler)
	//updateProjectHandler
	// http.HandleFunc("/EditAudition/", changeAuditionHandler)
	//END
	
	http.HandleFunc("/auditions/create", createRoleHandler)
	http.HandleFunc("/login", googleLoginHandler)
	http.HandleFunc("/GoogleCallback", googleCallbackHandler)
	http.HandleFunc("/auditions/browse/", roleBrowseHandler)
	http.HandleFunc("/auditions/", roleHandler)
	//http.HandleFunc("/upload/", uploadTestHandler)
	http.HandleFunc("/logout/", logoutHandler)
	//
	//update handlers
	http.HandleFunc("/api/v1/updateUser/", updateUserHandler)
	http.HandleFunc("/api/v1/submitRole/", submitRoleHandler)
	http.HandleFunc("/api/v1/submitAudition/", submitAuditionHandler)
	http.HandleFunc("/api/v1/submitRoleComment/", submitRoleCommentHandler)
	http.HandleFunc("/api/v1/submitAuditionComment/", submitAuditionCommentHandler)
	http.HandleFunc("/api/v1/submitProject/", submitProjectHandler)
	http.HandleFunc("/api/v1/getRole/", getRoleHandler)
	http.HandleFunc("/api/v1/submitTeam/", submitTeamHandler)
	http.HandleFunc("/api/v1/submitContest/", submitContestHandler)
	//Listen on port 80
	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}