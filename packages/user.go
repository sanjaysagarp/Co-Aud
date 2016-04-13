package user

import (
	"log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//User struct
type User struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	DisplayName string
	Username string
	AboutMe string
	PersonalWebsite string
	FacebookURL string
	InstagramURL string
	TwitterURL string
	VideoSubmission []Work
	ContestTeamNames []string
	Email string
	
}

//Work struct defines a person's personal work
type Work struct {
	Title string
	URL string
	ShortDescription string
	Description string
	Cast []Cast
	PostedDate string
	PostedTime string
}

//Cast struct
type Cast struct {
	Username string
	Role string
}

//Contest struct
type Contest struct {
	CreatedBy string
	ShortDescription string
	Description string
	ParticipatingTeams []Team
	StartDate string
	EndDate string
}

//Team struct
type Team struct {
	UserNames []Cast
	TeamName string
}

//Role struct
type Role struct {
	Username string
	ShortDescription string
	Description string
	VotesUp int
	VotesDown int
}

//Comment struct
type Comment struct {
	Username string
	Message string
	Replies []Comment
}

//NewUser creates a new user after signed in with google
func NewUser(email string, displayName string) *User{
	return &User{Email: email, DisplayName: displayName}
}

//FindUser searches for the user
func FindUser(email string) *User {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")
	
	result := &User{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		fmt.Println("User not found, creating one")
		return nil
	}
	return result
}


//InsertUser adds the user to the db
func InsertUser(user *User) {
	session, err := mgo.Dial("127.0.0.1:27018")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")
	err = c.Insert(&User{Email: user.Email, DisplayName: user.DisplayName})
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.DisplayName + " added with email " + user.Email)
}