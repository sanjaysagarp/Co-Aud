package work
// //127.0.0.1:27018
import (
	// "log"
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanjaysagarp/Co-Aud/packages/user"
	"strings"
)

//Cast struct
type Cast struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User *user.User
	Role string
}

//Work struct defines a person's personal work
type Work struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	URL string
	ShortDescription string
	Description string
	Cast []Cast // [] Return to the original later **
	PostedDate time.Time
	User *user.User
}

func (w *Work) GetYoutubeID() string {
	youtubeCode := strings.Split(w.URL, "=")
    return youtubeCode[1]
}

//NewWork creates a new instance of work
func NewWork(title string, url string, shortDescription string, description string, cast []Cast, user *user.User, id bson.ObjectId) *Work {
	return &Work{Title: title, URL: url, ShortDescription: shortDescription, Description : description, Cast: cast, PostedDate: time.Now(), User : user, Id: id}
}

//NewCast creates a new instance of cast
func NewCast(user *user.User, role string) Cast {
	return Cast{User: user, Role: role}
}

//InsertWork inserts a work into the works collection
func InsertWork(work *Work) {
	session, err := mgo.Dial("127.0.0.1")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	err = c.Insert(&Work{Title: work.Title, URL: work.URL, ShortDescription: work.ShortDescription, Description: work.Description, Cast: work.Cast, PostedDate: work.PostedDate, User: work.User, Id: work.Id})
	
	if err != nil {
		panic(err)
	}
}

//InsertCast inserts a new cast into a work
func InsertCast(cast *Cast) {
	session, err := mgo.Dial("127.0.0.1")
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

//FindCast finds casting for work
func FindCast(work *Work) []Cast{
	session, err := mgo.Dial("127.0.0.1")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	
	result := []Cast{}
	err = c.Find(bson.M{"Id": work.Id}).All(&result)
	if err != nil {
		fmt.Println("Work now found")
		return nil
	}
	return result
}

//FindWorks finds works for all selected
func FindWork(id string) *Work{
	session, err := mgo.Dial("127.0.0.1")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	
	result := &Work{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		fmt.Println("Work now found")
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
func FindWorks(q interface{}, skip int, limit int) ([]Work, int) {
	session, err := mgo.Dial("127.0.0.1")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	result := []Work{}
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