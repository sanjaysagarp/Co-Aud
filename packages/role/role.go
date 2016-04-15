package role

import (
	// "log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/sanjaysagarp/Co-Aud/packages/work"
)

//Contest struct
type Contest struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	CreatedBy string
	ShortDescription string
	Description string
	ParticipatingTeams []Team
	StartDate string
	EndDate string
}

//Team struct
type Team struct {
	UserNames []work.Cast
	TeamName string
}

//Role struct - posting
type Role struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Title string
	Username string
	ShortDescription string
	Description string
	Comment []Comment
	VotesUp int
	VotesDown int
}

//Comment struct - Maybe include audio clip
type Comment struct {
	Username string
	Message string
	TimeStamp string
	Replies []Comment
}

//NewComment creates an instance of a new comment and returns it
//TODO: FILL OUT FIELDS
func NewComment() *Comment {
	return &Comment{}
}

//NewRole creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewRole() *Role {
	return &Role{}
}

//NewTeam creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewTeam() *Team {
	return &Team{}
}

//NewContest creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS 
func NewContest() *Contest {
	return &Contest{}
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
	
	
	
}

//InsertContest inserts contest into db
func InsertContest(contest *Contest) {
	
}

//InsertRole inserts role into db
func InsertRole(role *Role) {
	
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