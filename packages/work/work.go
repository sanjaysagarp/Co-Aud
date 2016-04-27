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
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID string
	Role string
}

//Work struct defines a person's personal work
type Work struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	URL string
	ShortDescription string
	Description string
	Cast []Cast
	PostedDate time.Time
	UserEmail string
}

//NewWork creates a new instance of work
func NewWork(title string, url string, shortDescription string, description string, cast []Cast, userEmail string) *Work {
	return &Work{
		Title: title,
		URL: url,
		ShortDescription: shortDescription,
		Cast: cast,
		PostedDate: time.Now(),
		UserEmail: userEmail
	}
}

//NewCast creates a new instance of cast
func NewCast(userID string, role string) *Cast {
	return &Cast{UserID: userID, Role: role}
}

//InsertWork inserts a work into the works collection
func InsertWork(work *Work) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	err = c.Insert(&Work{Title: work.Title, URL: work.URL, ShortDescription: work.ShortDescription, Cast: work.Cast,
												PostedDate: work.PostedDate, PostedTime, work.PostedTime, UserEmail: work.UserEmail})
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(
		"Title: " + work.Title +
		"\nURL: " + work.URL +
		"\nShort Description: " + work.ShortDescription +
		"\nCast: " + work.Cast +
		"\nPosted Date: " + work.PostedDate +
		"\nPosted Time: " + work.PostedTime +
		"\nUserEmail: " + work.UserEmail)
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
	err = c.Find(bson.M{"Id": bson.ObjectIdHex(work.Id)}).One(&result).Cast //<=============this should get the cast array from works
	if err != nil {
		fmt.Println("Work now found")
		return nil
	}
	return result
}

//FindWorks finds works for all selected
func FindWorks(userEmail string) []Work{
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	
	result := &Work[]
	err = c.Find(bson.M{"UserEmail": userEmail})
	if err != nil {
		fmt.Println("Work now found")
		return nil
	}
	return result
}

//FindWorks finds work based on string and returns slice of Work
func FindWorks(title string) []Work{
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("works")
	
	result := &Work[]
	err = c.Find(bson.M{"Title": title})
	if err != nil {
		fmt.Println("Work now found")
		return nil
	}
	return result
}
