package role

import (
	// "log"
	// "fmt"
	// "gopkg.in/mgo.v2"
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