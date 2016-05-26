package user

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//User struct
type User struct {
	Id                bson.ObjectId `json:"id" bson:"_id,omitempty"`
	DisplayName       string
	Email             string
	Title             string
	AboutMe           string
	PersonalWebsite   string
	FacebookURL       string
	InstagramURL      string
	TwitterURL        string
	ContestTeamNames  []string
	ProfilePictureURL string
	AwsPictureURL     string
	JoinDate          time.Time
}

//NewUser creates a new user after signed in with google
func NewUser(email string, displayName string) *User {
	return &User{Email: email, DisplayName: displayName, ProfilePictureURL: "/public/img/default_profile_pic.png", JoinDate: time.Now()}
}

//NewChangeUser creates a new user with most fields
func NewChangeUser(displayName string, title string, aboutMe string, personalWebsite string, facebookURL string, instagramURL string, twitterURL string) *User {

	return &User{DisplayName: displayName, Title: title, PersonalWebsite: personalWebsite, AboutMe: aboutMe, FacebookURL: facebookURL, TwitterURL: twitterURL, InstagramURL: instagramURL}
}

//FindUser searches for the user
func FindUser(email string) *User {
	session, err := mgo.Dial("127.0.0.1")
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

//FindUser searches for the user
func FindUserById(id string) *User {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")

	result := &User{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		fmt.Println("User not found")
		return nil
	}
	return result
}

//InsertUser adds the user to the db
func InsertUser(user *User) {
	session, err := mgo.Dial("127.0.0.1")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")
	err = c.Insert(user)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.DisplayName + " added with email " + user.Email)
}

//UpdateUser updates a user with the given id and handler made struct
func UpdateUser(id string, user *User) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")

	//this shit is erasing the fields? Need to check for consistency
	change := bson.M{
		"$set": bson.M{
			"displayname":     user.DisplayName,
			"title":           user.Title,
			"aboutme":         user.AboutMe,
			"personalwebsite": user.PersonalWebsite,
			"facebookurl":     user.FacebookURL,
			"instagramurl":    user.InstagramURL,
			"twitterurl":      user.TwitterURL}}
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
	if err != nil {
		panic(err)
	}
}

//UpdateUserPicture updates a user's profile picture -- NEED TO TEST
func UpdateUserPicture(id string, URL string, AWSURL string, user *User) {

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("users")

	change := bson.M{"$set": bson.M{
		"profilepictureurl": URL,
		"awspictureurl":     AWSURL}}
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
	if err != nil {
		panic(err)
	}
}
