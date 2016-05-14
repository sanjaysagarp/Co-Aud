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
	ParticipatingTeams []*mgo.DBRef
	ImageUrl string
	StartDate time.Time
	EndDate time.Time
}

//Team struct
type Team struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Users []*mgo.DBRef
	TeamName string
	Motto string
}

//Role struct - posting
type Role struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title string
	ImageUrl string
	User *mgo.DBRef
	Traits []string
	Description string
	Script string
	Gender string
	Age int
	TimeStamp time.Time
	Deadline time.Time
	Comment []*mgo.DBRef
	Audition []*mgo.DBRef
}

func (r *Role) GetComments() []*Comment {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	result := []*Comment{}
	for i, c := range r.Comment {
		oneResult := &Comment{}
		err = session.FindRef(c).One(oneResult)
		if err != nil {
			fmt.Println("error happened at index: ", i)
			panic(err)
		}
		result = append(result, oneResult)
	}
    return result
}

func (r *Role) GetUser() *user.User {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	result := &user.User{}
	err = session.FindRef(r.User).One(result)
	if err != nil {
		panic(err)
	}
    return result
}

func (r *Role) GetAuditions() []*Audition {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	result := []*Audition{}
	for i, a := range r.Audition {
		oneResult := &Audition{}
		err = session.FindRef(a).One(oneResult)
		if err != nil {
			fmt.Println("error happened at index: ", i)
			panic(err)
		}
		result = append(result, oneResult)
	}
    return result
}

//Comment struct - Maybe include audio clip
type Audition struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User *mgo.DBRef
	AttachmentUrl string
	TimeStamp time.Time
	Comment []*mgo.DBRef
}

func (a *Audition) GetComments() []*Comment {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	result := []*Comment{}
	for i, c := range a.Comment {
		oneResult := &Comment{}
		err = session.FindRef(c).One(oneResult)
		if err != nil {
			fmt.Println("error happened at index: ", i)
			panic(err)
		}
		result = append(result, oneResult)
	}
    return result
}

func (a *Audition) GetUser() *user.User {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	result := &user.User{}
	err = session.FindRef(a.User).One(result)
	if err != nil {
		panic(err)
	}
    return result
}

//Comment struct - Maybe include audio clip
type Comment struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	User *mgo.DBRef
	Message string
	TimeStamp time.Time
}

func (c *Comment) GetUser() *user.User {
	session, err := mgo.Dial("127.0.0.1:27018")
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	result := &user.User{}
	err = session.FindRef(c.User).One(result)
	if err != nil {
		panic(err)
	}
    return result
}

//NewComment creates an instance of a new comment and returns it
//TODO: FILL OUT FIELDS
func NewComment(user *user.User, message string, id bson.ObjectId) *Comment {
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	return &Comment{User : dbRefUser, Message : message, TimeStamp : time.Now(), Id: id}
}

//NewRole creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func NewRole(title string, user *user.User, description string, script string, deadline time.Time, traits []string, age int, gender string, id bson.ObjectId, imageUrl string) *Role {
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	return &Role{Title: title, User: dbRefUser, Description: description, Script: script, TimeStamp: time.Now(), Deadline: deadline, Traits: traits, Age: age, Gender: gender, ImageUrl: imageUrl, Id: id}
}

//NewTeam creates an instance of a new role and returns it
//TODO: FILL OUT FIELDS
func InsertNewTeam(users []*user.User, teamName string, motto string) *Team {
	session, err := mgo.Dial("127.0.0.1:27018")
	
	//session, err := mgo.Dial("127.0.0.1")
	if err != nil {
			panic(err)
	}
	defer session.Close()
	var dbRefUsers []*mgo.DBRef
	for _,user := range users {
		dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
		dbRefUsers = append(dbRefUsers, dbRefUser)
	}
	//dbRefContest := &mgo.DBRef{Collection: "contests", Id: contest.id, Database: "CoAud"}
	team := &Team{Users: dbRefUsers, TeamName: teamName, Motto: motto}
	
	
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("teams")
	err = c.Insert(team)
	
	if err != nil {
		panic(err)
	}
	return team

}

func NewAudition(user *user.User, attachmentUrl string, id bson.ObjectId) *Audition {
	dbRefUser := &mgo.DBRef{Collection: "users", Id: user.Id, Database: "CoAud"}
	return &Audition{User: dbRefUser, AttachmentUrl: attachmentUrl, TimeStamp: time.Now(), Id: id}
}

//NewContest creates an instance of a new role and returns it
func NewContest(title string, description string, imageUrl string, endDate time.Time) *Contest {
	return &Contest{Title: title,Description: description,ImageUrl: imageUrl,StartDate: time.Now(),EndDate: endDate}
}

//InsertComment takes a role and inserts a comment into the role's comment array
//Need to grab data in handler and create a new comment struct
//TODO: insert comment to role (db)

//GENERIC INSERT COMMENT
//
func InsertComment(comment *Comment, collection string, id string, recentOrder bool) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("comments")

	err = c.Insert(comment)
	if err != nil {
		panic(err)
	}
	
	c = session.DB("CoAud").C(collection)
	
	dbRefComment := &mgo.DBRef{Collection: "comments", Id: comment.Id, Database: "CoAud"}

	// store comment into slice so we can push it to the top in mongo
	// $push with $position requires $each which requires an array/slice
	var newComment []*mgo.DBRef
	newComment = append(newComment, dbRefComment)
	if recentOrder {
		change := bson.M{"$push": bson.M{"comment": bson.M{"$each": newComment, "$position": 0}}}
		err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
	} else {
		change := bson.M{"$push": bson.M{"comment": bson.M{"$each": newComment}}}
		err = c.Update(bson.M{"_id": bson.ObjectIdHex(id)}, change)
	}
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
	
	//insert new audition
	c := session.DB("CoAud").C("auditions")

	err = c.Insert(audition)
	if err != nil {
		panic(err)
	}
	fmt.Println("audition.Id: ")
	fmt.Println(audition.Id)
	
	//reference audition to role
	c = session.DB("CoAud").C("roles")
	
	dbRefAudition := &mgo.DBRef{Collection: "auditions", Id: audition.Id, Database: "CoAud"}
	
	change := bson.M{"$push": bson.M{"audition": dbRefAudition}}
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

func (contest *Contest) InsertTeam(team *Team) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("contests")
	// Find contest, then insert into contest by contest name
	
	//contest.ParticipatingTeams = append(contest.ParticipatingTeams, team)
	//box.AddItem(item1)
	//change := bson.M{"$set": bson.M{"ParticipatingTeams": contest.AddItem(team)}}
	
	dbRefTeam:= &mgo.DBRef{Collection: "teams", Id: team.Id, Database: "CoAud"}
	change := bson.M{"$push": bson.M{"ParticipatingTeams": bson.M{"$each": dbRefTeam}}}
	err = c.Update(bson.M{"_id": contest.Id}, change)
	
	if err != nil {
		panic(err)
	}
}

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