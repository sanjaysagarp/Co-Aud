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
func NewRole(title string, userEmail string, description string, script string, traits []string) *Role {
	return &Role{Title: title, UserEmail: userEmail, Description: description, Script: script, Traits: traits}
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
func InsertComment(role *Role, comment Comment) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	role.Comment = append(role.Comment, comment)
	change := bson.M{"$set": bson.M{"Comment": role.Comment}}
	//err = c.Insert(&Role{User: cast.User, Role: cast.Role})
	
	err = c.Update(bson.M{"_id": role.Id}, change)
	//update role - add comment to slice
}

func InsertComment1(audition *Audition, comment Comment) {
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
	audition.Comment = append(audition.Comment, comment)
	change := bson.M{"$set": bson.M{"Comment": audition.Comment}}
	err = c.Update(bson.M{"_id": audition.Id}, change)
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
	c := session.DB("CoAud").C("contests")

	err = c.Insert(&Contest{Title: contest.Title, Description: contest.Description, ImageUrl: contest.ImageUrl, StartDate: contest.StartDate, EndDate: contest.EndDate})
	if err != nil {
		panic(err)
	}
}

//InsertContest inserts contest into db
func InsertAudition(audition *Audition) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")

	err = c.Insert(&Audition{UserEmail: audition.UserEmail, AttatchmentUrl: audition.AttatchmentUrl, TimeStamp: audition.TimeStamp, Comment: audition.Comment})
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
	c := session.DB("CoAud").C("roles")

	err = c.Insert(&Role{Title: role.Title,UserEmail: role.UserEmail,Description: role.Description,Script: role.Script,TimeStamp: role.TimeStamp,Deadline: role.Deadline,Traits: role.Traits})
	if err != nil {
		panic(err)
	}
}

//InsertTeam insert team into db within a Contest
func InsertTeam(contest *Contest, team Team) {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("contests")
	// Find contest, then insert into contest by contest name
	
	// bloop := contest.ParticipatingTeams
	// bloop = append(bloop, team)
	contest.ParticipatingTeams = append(contest.ParticipatingTeams, team)
	//box.AddItem(item1)
	//bloop := make([]string, 0)
	//teamStar := NewTeam(bloop, "teamName", "21321")
	change := bson.M{"$set": bson.M{"ParticipatingTeams": contest.ParticipatingTeams}}
	err = c.Update(bson.M{"_id": contest.Id}, change)
	
	//err = c.Insert(&Team{UserNames: team.Username,TeamName: team.teamName})
	
	if err != nil {
		panic(err)
	}
}

func (contest *Contest) AddItem(team Team) []Team {
    contest.ParticipatingTeams = append(contest.ParticipatingTeams, team)
    return contest.ParticipatingTeams
}

//FindRoles searches for all roles
//TODO: query db for roles and add to result, then return roles
func FindRoles(role *Role) []Role {
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

//FindRole searches for the selected role
//TODO: query db for roles and add to result, then return roles
func FindSearchedRole(title string) Role {
	session, err := mgo.Dial("127.0.0.1:27018")
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("CoAud").C("roles")
	
	result := Role{}
	err = c.Find(bson.M{"Title": title}).One(&result)
	if err != nil {
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