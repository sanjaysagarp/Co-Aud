package work

import (
	// "log"
	// "fmt"
	// "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Cast struct
type Cast struct {
	Username string
	Role string
}

//Work struct defines a person's personal work
type Work struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Title string
	URL string
	ShortDescription string
	Description string
	Cast []Cast
	PostedDate string
	PostedTime string
}

