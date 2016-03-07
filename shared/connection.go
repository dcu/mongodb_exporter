package shared

import (
	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

const (
	dialMongodbTimeout = 10 * time.Second
	syncMongodbTimeout = 1 * time.Minute
)

func MongoSession(uri string) *mgo.Session {
	dialInfo, err := mgo.ParseURL(uri)
	if err != nil {
		glog.Errorf("Cannot connect to server using url %s: %s", uri, err)
		return nil
	}

	dialInfo.Direct = true // Force direct connection
	dialInfo.Timeout = dialMongodbTimeout

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		glog.Errorf("Cannot connect to server using url %s: %s", uri, err)
		return nil
	}
	session.SetMode(mgo.Eventual, true)
	session.SetSyncTimeout(syncMongodbTimeout)
	session.SetSocketTimeout(0)
	return session
}
