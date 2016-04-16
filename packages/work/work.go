package work

import (
	// "log"
	"fmt"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanjaysagarp/Co-Aud/packages/user"
)

//Cast struct
type Cast struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	User bson.ObjectId
	Role string
}

//Work struct defines a person's personal work
type Work struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Title string
	URL string
	ShortDescription string
	Description string
	Cast []Cast
	PostedDate string
	PostedTime string
	User bson.ObjectId
}

//NewWork creates a new instance of work
func NewWork(title string, url string, shortDescription string, description string, cast []Cast, User bson.ObjectId) *Work {
	return &Work{
		Title: title,
		URL: url,
		ShortDescription: shortDescription,
		Cast: cast,
		PostedDate: time.Now().Date(),
		PostedTime: time.Now().Clock(),
		User: user
	}
}

//NewCast creates a new instance of cast
func NewCast(user bson.ObjectId, role string) *Cast {
	return &Cast{User: user, Role: role}
}

//InsertWork inserts a work into the works collection
//----CAUTION----
//Dont worry about this meow -- need to get user session for this
func InsertWork(work *Work) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")
	
	
}

//InsertCast inserts a new cast into a work
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
		log.Fatal(err)
	}
	fmt.Println("User with ID: " + cast.User + " added as " + cast.Role)
}

//FindCast finds casting for work
func FindCast(work *Work) []Cast{
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	
	result := &Cast[]
	err = c.Find(bson.M{"ID": work.ID}).One(&result).Cast //<=============this should get the cast array from works
	if err != nil {
		fmt.Println("Work now found")
		return nil
	}
	return result
}

//FindWorks finds works for all selected
//----CAUTION----
//Dont worry about this meow -- need to get user session for this
func FindWorks(user *User) []Work{
	
}

//FindWorks finds work based on string and returns slice of Work
func FindWorks(title string) []Work{
	
}
