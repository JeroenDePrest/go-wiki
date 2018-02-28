package mongoDB

import (
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func Create() {
	var err error
	session, err = mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
}

func Session() *mgo.Session {
	return session
}
