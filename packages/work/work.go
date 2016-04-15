package work

import (
	// "log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanjaysagarp/Co-Aud/packages/user"
)

//Cast struct
type Cast struct {
	Username string
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
	Author user.User
}

//NewWork creates a new instance of work
func NewWork() *Work {
	return &Work{}
}

//NewCast creates a new instance of cast
func NewCast() *Cast {
	return &Cast{}
}

//InsertWork inserts a work into a users acct
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
	
}

//FindCast finds casting for work
func FindCast(work *Work) []Cast{
	
}

//FindWorks finds works for all selected
//----CAUTION----
//Dont worry about this meow -- need to get user session for this
func FindWorks(author *User) []Work{
	
}

//FindWorks finds work based on string and returns slice of Work
func FindWorks(title string) []Work{
	
}
