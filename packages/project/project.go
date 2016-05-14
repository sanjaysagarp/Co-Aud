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
	User *user.User
	Role string
}

//Project struct defines a person's personal project
type Project struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	URL string
	ShortDescription string
	Description string
	Cast []Cast // [] Return to the original later **
	PostedDate time.Time
	User *user.User
}

func (w *Project) GetYoutubeID() string {
    r, _ := regexp.Compile(`^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#\&\?]*).*`)
    return r.FindAllStringSubmatch(w.URL, -1)[0][7]
}

//NewProject creates a new instance of project
func NewProject(title string, url string, shortDescription string, description string, cast []Cast, user *user.User, id bson.ObjectId) *Project {
	return &Project{Title: title, URL: url, ShortDescription: shortDescription, Description : description, Cast: cast, PostedDate: time.Now(), User : user, Id: id}
}

//NewCast creates a new instance of cast
func NewCast(user *user.User, role string) Cast {
	return Cast{User: user, Role: role}
}

//InsertProject inserts a project into the projects collection
func InsertProject(project *Project) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("projects")
	err = c.Insert(&Project{Title: project.Title, URL: project.URL, ShortDescription: project.ShortDescription, Description: project.Description, Cast: project.Cast, PostedDate: project.PostedDate, User: project.User, Id: project.Id})
	
	if err != nil {
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
	err = c.Insert(&Cast{User: cast.User, Role: cast.Role})
	
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