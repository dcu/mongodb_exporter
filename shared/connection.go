package shared

import (
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

func MongoSession(uri string) *mgo.Session {
	session, err := mgo.Dial(uri)
	if err != nil {
		glog.Errorf("Cannot connect to server using url: %s", uri)
		return nil
	}
	session.SetMode(mgo.Eventual, true)
	session.SetSocketTimeout(0)
	return session
}
