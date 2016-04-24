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
	ContestTeamNames []string
	Email string
	
}

//NewUser creates a new user after signed in with google
func NewUser(email string, displayName string) *User{
	return &User{Email: email, DisplayName: displayName}
}

//FindUser searches for the user
func FindUser(email string) *User {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	//fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")
	
	result := &User{}
	err = c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		fmt.Println("User not found")
		return nil
	}
	return result
}

//InsertUser adds the user to the db
func InsertUser(user *User) {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
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