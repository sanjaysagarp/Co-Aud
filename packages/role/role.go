package role

import (
	// "log"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/sanjaysagarp/Co-Aud/packages/user"
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
	ContestId string
}

//Role struct - posting
type Role struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	UserEmail string
	Traits []string
	Description string
	Script string
	Gender string
	Age int
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
	return &Comment{UserEmail : userEmail, Message : message, TimeStamp : time.Now()}
}

//NewRole creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS

		//TimeStamp: time.Now(),
		//Deadline: deadline,
func NewRole(title string, userEmail string, description string, script string, deadline time.Time, traits []string, age int, gender string, id bson.ObjectId) *Role {
	return &Role{Title: title, UserEmail: userEmail, Description: description, Script: script, TimeStamp: time.Now(), Deadline: deadline, Traits: traits, Age: age, Gender: gender, Id: id}
}

//NewTeam creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewTeam(userNames []string, teamName string, contestId string) *Team {
	session, err := mgo.Dial("127.0.0.1:27018")
	team := &Team{UserNames: userNames, TeamName: teamName, ContestId: contestId}
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("teams")
	err = c.Insert(team)
	
	if err != nil {
		panic(err)
	}
	return team

}

func NewAudition(userEmail string, attachmentUrl string) *Audition {
	return &Audition{UserEmail: userEmail, AttachmentUrl: attachmentUrl, TimeStamp: time.Now()}
}

//NewContest creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS 
func NewContest(title string, description string, imageUrl string, endDate time.Time) *Contest {
	return &Contest{Title: title,Description: description,ImageUrl: imageUrl,StartDate: time.Now(),EndDate: endDate}
}

//InsertComment takes a role and inserts a comment into the role's comment array
//Need to grab data in handler and create a new comment struct
//TODO: insert comment to role (db)
// func InsertComment(role *Role, comment *Comment) {
// 	session, err := mgo.Dial("127.0.0.1:27018")
// 	fmt.Println("connected")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("CoAud").C("role")
	

// 	change := bson.M{"$set": bson.M{"Comment": role.Comment.append(comment)}}
// 	err = c.Insert(&Role{User: cast.User, Role: cast.Role})
	
// 	err = c.Update(bson.M{"_id": bson.ObjectIdHex(role.Id)}, change)
// 	//update role - add comment to slice
// }

// func InsertComment1(audition *Audition, comment *Comment) {
// 	session, err := mgo.Dial("127.0.0.1:27018")
// 	fmt.Println("connected")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("CoAud").C("audition")

// 	//find audition
// 	//update audition - add comment to slice	
// 	change := bson.M{"$set": bson.M{"Comment": audition.Comment.append(comment)}}
// 	err = c.Update(bson.M{"_id": bson.ObjectIdHex(role.Id)}, change)
// }

//InsertContest inserts contest into db
func InsertContest(contest *Contest) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("contests")

	err = c.Insert(&Contest{Title: contest.Title, Description: contest.Description, ImageUrl: contest.ImageUrl, StartDate: contest.StartDate, EndDate: contest.EndDate})
	if err != nil {
		panic(err)
	}
}

//InsertRole inserts role into db
func InsertRole(role *Role) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		fmt.Println("InsertRole IS THE PROBLEM")
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")

	err = c.Insert(&Role{Title: role.Title,UserEmail: role.UserEmail,Description: role.Description,Script: role.Script,TimeStamp: role.TimeStamp,Deadline: role.Deadline,Traits: role.Traits, Gender: role.Gender, Age: role.Age, Id: role.Id})
	if err != nil {
		panic(err)
	}
}

//InsertTeam insert team into db within a Contest
// func InsertTeam(contest *Contest, team *Team) {
// 	session, err := mgo.Dial("127.0.0.1:27018")
// 	fmt.Println("connected")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("CoAud").C("contest")
// 	// Find contest, then insert into contest by contest name
	
// 	//contest.ParticipatingTeams = append(contest.ParticipatingTeams, team)
// 	//box.AddItem(item1)
// 	change := bson.M{"$set": bson.M{"ParticipatingTeams": contest.AddItem(team)}}
// 	err = c.Update(bson.M{"_id": bson.ObjectIdHex(contest.Id)}, change)

// 	err = c.Insert(&Team{UserNames: team.Username,TeamName: team.teamName})
	
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func (contest *Contest) AddItem(team Team) []Team {
//     contest.ParticipatingTeams = append(contest.ParticipatingTeams, team)
//     return contest.ParticipatingTeams
// }

// FindRoles searches for all roles
// TODO: query db for roles and add to result, then return roles
func FindRoles() []Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	result := []Role{}
	err = c.Find(nil).All(&result)
	if err != nil {
		panic(err)
	}
	return result
}

func FindRolesByUserEmail(userEmail string) []Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	result := []Role{}
	err = c.Find(bson.M{"useremail": userEmail}).All(&result)
	if err != nil {
		panic(err)
	}
	return result
}

//FindRole searches for the selected role
//TODO: query db for roles and add to result, then return roles
func FindRole(id string) *Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	result := &Role{}
	err = c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		fmt.Println("Role not found")
		panic(err)
	}
	return result
}

//FindContests searches for all contests
//TODO: query db for contests and add to result, then return contests
func FindContests(title string) []Contest {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("contest")
	
	result := []Contest{}
	err = c.Find(nil).All(&result)
	if err != nil {
		panic(err)
	}
	
	return result
}

//FindContest searches for the user
//TODO: query db for a single contest and add to result, then return roles
func FindSearchedContest(title string) []Contest {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("contest")
	result := []Contest{}
	err = c.Find(bson.M{"Title": title}).One(&result)
	
	
	return result
}