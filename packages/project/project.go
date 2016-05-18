package project
// //127.0.0.1:27018
import (
	// "log"
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanjaysagarp/Co-Aud/packages/user"
	// "strings"
	"regexp"
)

//Cast struct
type Cast struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User *mgo.DBRef
	Role string
}

//Project struct defines a person's personal project
type Project struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string 
	URL string
	ShortDescription string
	Description string
	Cast []*mgo.DBRef
	PostedDate time.Time
	User *mgo.DBRef
}

func (w *Project) GetYoutubeID() string {
    r, _ := regexp.Compile(`^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#\&\?]*).*`)
    return r.FindAllStringSubmatch(w.URL, -1)[0][7]
}

//NewProject creates a new instance of project
func NewProject(title string, url string, shortDescription string, description string, casts []*Cast, user *user.User, id bson.ObjectId) *Project {
	var dbRefCasts []*mgo.DBRef
	for _, cast := range casts {
		fmt.Println("Cast ID: " + cast.Id)
		dbRefCast := &mgo.DBRef{Collection: "casts", Id: cast.Id, Database: "CoAud"}
		dbRefCasts = append(dbRefCasts, dbRefCast)
	}
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	fmt.Println(dbRefUser)
	fmt.Println("project id: ", id)
	return &Project{Id: id, Title: title, URL: url, ShortDescription: shortDescription, Description: description, Cast: dbRefCasts, PostedDate: time.Now(), User: dbRefUser}
}

//NewCast creates a new instance of cast
func NewCast(user *user.User, role string, id bson.ObjectId) *Cast {
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	return &Cast{Id: id, User: dbRefUser, Role: role}
}

//InsertProject inserts a project into the projects collection
func InsertProject(project *Project) {
	session, err := mgo.Dial("127.0.0.1:27018")
	if err != nil {
		fmt.Println("not connected")
		panic(err)
	}
	fmt.Println("connected")

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("projects")
	err = c.Insert(project)
	if err != nil {
		fmt.Println("insert project fails")
		panic(err)
	}
}

//InsertCast inserts a new cast into a project
func InsertCast(cast *Cast) {
	session, err := mgo.Dial("127.0.0.1:27018")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("casts")
	err = c.Insert(cast)
	
	if err != nil {
		panic(err)
	}
}

//FindCast finds casting for project
func FindCast(project *Project) []Cast{
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("projects")
	
	result := []Cast{}
	err = c.Find(bson.M{"Id": project.Id}).All(&result)
	if err != nil {
		fmt.Println("Work now found")
		return nil
	}
	return result
}

//FindProjects finds projects for all selected
func FindProject(id string) *Project{
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("projects")
	
	result := &Project{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		fmt.Println("Project now found")
		panic(err)
	}
	return result
}

// //FindWorks finds work based on string and returns slice of Work
// func FindWorks2(title string) Work{
// 	session, err := mgo.Dial("127.0.0.1:27018")
// 	fmt.Println("connected")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("CoAud").C("works")
	
// 	result := Work{}
// 	err = c.Find(bson.M{"Title": title}).One(&result)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return result
// }


// FindRoles searches for all roles
// Optional param: q = nil, skip = 0, limit = -1
func FindProjects(q interface{}, skip int, limit int) ([]Project, int) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("projects")
	result := []Project{}
	err = c.Find(q).Skip(skip).Limit(limit).Sort("-posteddate").All(&result)
	if err != nil {
		panic(err)
	}
	resultCount, err := c.Count()
	if err != nil {
		panic(err)
	}
	
	return result, resultCount
}