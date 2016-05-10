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
	Users []*user.User
	TeamName string
	ContestId string
}

//Role struct - posting
type Role struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	ImageUrl string
	User *user.User
	Traits []string
	Description string
	Script string
	Gender string
	Age int
	TimeStamp time.Time
	Deadline time.Time
	Comment []*Comment
	Audition []Audition
}

//Comment struct - Maybe include audio clip
type Audition struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User *user.User
	AttachmentUrl string
	TimeStamp time.Time
	Comment []*Comment
}

//Comment struct - Maybe include audio clip
type Comment struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User *user.User
	Message string
	TimeStamp time.Time
}

//NewComment creates an instance of a new comment and returns it
//TODO: FILL OUT FIELDS
func NewComment(user *user.User, message string) *Comment {
	return &Comment{User : user, Message : message, TimeStamp : time.Now()}
}

//NewRole creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewRole(title string, user *user.User, description string, script string, deadline time.Time, traits []string, age int, gender string, id bson.ObjectId, imageUrl string,) *Role {
	return &Role{Title: title, User: user, Description: description, Script: script, TimeStamp: time.Now(), Deadline: deadline, Traits: traits, Age: age, Gender: gender, ImageUrl: imageUrl, Id: id}
}

//NewTeam creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewTeam(users []*user.User, teamName string, contestId string) *Team {
	session, err := mgo.Dial("127.0.0.1:27018")
	team := &Team{Users: users, TeamName: teamName, ContestId: contestId}
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

func NewAudition(user *user.User, attachmentUrl string) *Audition {
	return &Audition{User: user, AttachmentUrl: attachmentUrl, TimeStamp: time.Now()}
}

//NewContest creates an instance of a new role and returns it
func NewContest(title string, description string, imageUrl string, endDate time.Time) *Contest {
	return &Contest{Title: title,Description: description,ImageUrl: imageUrl,StartDate: time.Now(),EndDate: endDate}
}

//InsertComment takes a role and inserts a comment into the role's comment array
//Need to grab data in handler and create a new comment struct
//TODO: insert comment to role (db)

//GENERIC INSERT COMMENT
func InsertComment(commentList []*Comment, comment *Comment, collection string, id string) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("comments")

	err = c.Insert(Comment{User: comment.User, Message: comment.Message, TimeStamp: comment.TimeStamp})
	if err != nil {
		panic(err)
	}
	
	c = session.DB("CoAud").C(collection)
	
	// store comment into slice so we can push it to the top in mongo
	// $push with $position requires $each which requires an array/slice
	var newComment []*Comment
	newComment = append(newComment, comment)
	change := bson.M{"$push": bson.M{"comment": bson.M{"$each": newComment, "$position": 0}}}
	
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
}

//InsertAudition inserts audition into db
func InsertAudition(audition *Audition, role *Role) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	change := bson.M{"$push": bson.M{"audition": audition}}
	err = c.Update(bson.M{"_id": role.Id}, change)
	if err != nil {
		panic(err)
	}
}

// func InsertComment2(role *Role, comment Comment) {
// 	session, err := mgo.Dial("127.0.0.1:27018")
// 	fmt.Println("connected")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("CoAud").C("roles")
	
// 	role.Comment = append(role.Comment, comment)
// 	change := bson.M{"$set": bson.M{"Comment": role.Comment}}
// 	//err = c.Insert(&Role{User: cast.User, Role: cast.Role})
	
// 	err = c.Update(bson.M{"_id": role.Id}, change)
// 	//update role - add comment to slice
// }

// func InsertComment1(audition *Audition, comment Comment) {
// 	session, err := mgo.Dial("127.0.0.1:27018")
// 	fmt.Println("connected")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("CoAud").C("auditions")

// 	//find audition
// 	//update audition - add comment to slice	
// 	audition.Comment = append(audition.Comment, comment)
// 	change := bson.M{"$set": bson.M{"Comment": audition.Comment}}
// 	err = c.Update(bson.M{"_id": audition.Id}, change)
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

	err = c.Insert(role)
	if err != nil {
		panic(err)
	}
}

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

// 	err = c.Insert(&Team{Users: team.Users,TeamName: team.teamName})
	
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func (contest *Contest) AddItem(team Team) []Team {
//     contest.ParticipatingTeams = append(contest.ParticipatingTeams, team)
//     return contest.ParticipatingTeams
// }

// FindRoles searches for all roles
// Optional param: q = nil, skip = 0, limit = -1
func FindRoles(q interface{}, skip int, limit int) ([]Role, int) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	result := []Role{}
	err = c.Find(q).Skip(skip).Limit(limit).Sort("-timestamp").All(&result)
	if err != nil {
		panic(err)
	}
	resultCount, err := c.Count()
	if err != nil {
		panic(err)
	}
	
	return result, resultCount
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
	c := session.DB("CoAud").C("contests")
	
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
	c := session.DB("CoAud").C("contests")
	result := []Contest{}
	err = c.Find(bson.M{"Title": title}).One(&result)
	
	return result
}