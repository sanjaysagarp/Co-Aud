package role

import (
	// "log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanjaysagarp/Co-Aud/packages/user"
	"time"
)

//Contest struct
type Contest struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	Description string
	ParticipatingTeams []Team
	ImageUrl string
	StartDate time.Time
	EndDate time.Time
}

//Team struct
type Team struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserNames []string
	TeamName string
}

//Role struct - posting
type Role struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	UserEmail string
	Traits []string
	Description string
	Script string
	TimeStamp time.Time
	Deadline time.Time
	Comment []Comment
	Audition []Audition
}

//Comment struct - Maybe include audio clip
type Audition struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserEmail string
	AttachmentUrl string
	TimeStamp time.Time
	Comment []Comment
}

//Comment struct - Maybe include audio clip
type Comment struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserEmail string
	Message string
	TimeStamp time.Time
}

//NewComment creates an instance of a new comment and returns it
//TODO: FILL OUT FIELDS
func NewComment(userEmail string, message string) *Comment {
	return &Comment{
		UserEmail : userEmail,
		Message : message,
		TimeStamp : time.Now()
	}
}

//NewRole creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewRole(title string, userEmail string, description string, script string, deadline time, traits []string) *Role {
	return &Role{
		Title: title,
		UserName: userEmail,
		Description: description,
		Script: script,
		TimeStamp: time.Now(),
		Deadline: deadline,
		Traits: traits
	}
}

//NewTeam creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewTeam(userNames []string, teamName string) *Team {
	return &Team{
		UserNames: userNames,
		TeamName: teamName
	}
}

func NewAudition(userEmail string, attachmentUrl string) *Team {
	return &Team{
		UserEmail: userEmail,
		AttachmentUrl: attachmentUrl,
		TimeStamp: time.Now()
	}
}

//NewContest creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS 
func NewContest(title string, description string, imageUrl string, endDate time) *Contest {
	return &Contest{
		Title: title,
		Description: description,
		ImageUrl: imageUrl,
		StartDate: time.Now(),
		EndDate: endDate
	}
}

//InsertComment takes a role and inserts a comment into the role's comment array
//Need to grab data in handler and create a new comment struct
//TODO: insert comment to role (db)
func InsertComment(role *Role, comment *Comment) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("role")
	

	change := bson.M{"$set": bson.M{"Comment": role.Comment.append(comment)}}
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(role.Id)}, change)
	//update role - add comment to slice
}

func InsertComment(audition *Audition, comment *Comment) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("audition")

	//find audition
	//update audition - add comment to slice	
	change := bson.M{"$set": bson.M{"Comment": audition.Comment.append(comment)}}
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(role.Id)}, change)
}

//InsertContest inserts contest into db
func InsertContest(contest *Contest) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("contest")

	err = c.Insert(&Contest{
		Title: contest.Title
		Description: contest.Description
		ImageUrl: contest.ImageUrl
		StartDate: contest.StartDate
		EndDate: contest.EndDate
	})
	if err != nil {
		panic(err)
	}
}

//InsertRole inserts role into db
func InsertRole(role *Role) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("role")

	err = c.Insert(&Role{
		Title: role.Title,
		UserEmail: role.UserEmail,
		Description: role.Description,
		Script: role.Script,
		TimeStamp: role.TimeStamp,
		Deadline: role.Deadline,
		Traits: role.Traits
	})
	if err != nil {
		panic(err)
	}
	
}

//InsertTeam inserts role into db within a Contest
func InsertTeam(contest *Contest, role *Role) {
	
}

//FindRoles searches for all roles
//TODO: query db for roles and add to result, then return roles
func FindRoles(title string) []Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	result := []Role{}
	
	
	
	
	return result
}

//FindRole searches for the selected role
//TODO: query db for roles and add to result, then return roles
func FindRole(title string) Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	result := Role{}
	
	
	
	
	return result
}

//FindContests searches for all contests
//TODO: query db for contests and add to result, then return contests
func FindContests(title string) []Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	result := []Contest{}
	
	
	
	
	return result
}

//FindContest searches for the user
//TODO: query db for a single contest and add to result, then return roles
func FindContest(title string) []Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	result := Contest{}
	
	
	
	
	return result
}