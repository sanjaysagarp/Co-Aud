package user

import (
	"log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//User struct
type User struct {
	displayName string
	username string
	aboutMe string
	personalWebsite string
	facebookLink string
	instagramLink string
	twitterLink string
	videoSubmission []Work
	contestTeamNames []string
	email string
	
}

//Work struct defines a person's personal work
type Work struct {
	title string
	url string
	shortDescription string
	description string
	cast []Cast
	postedDate string
	postedTime string
}

//Cast struct
type Cast struct {
	username string
	role string
}

//Contest struct
type Contest struct {
	createdBy string
	shortDescription string
	description string
	participatingTeams []Team
	startDate string
	endDate string
}

//Team struct
type Team struct {
	userNames []Cast
	teamName string
}

//Role struct
type Role struct {
	username string
	shortDescription string
	description string
	votesUp int
	votesDown int
}

//Comment struct
type Comment struct {
	username string
	message string
	replies []Comment
}

//NewUser creates a new user after signed in with google
func NewUser(email string, displayName string) *User{
	return &User{email: email, displayName: displayName}
}

//FindUser searches for the user
func FindUser(email string) *User {
	session, err := mgo.Dial("127.0.0.1:27018")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")
	
	result := &User{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
			log.Fatal(err)
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

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("CoAud").C("users")
	err = c.Insert(user)
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.displayName + " added with email " + user.email)
}