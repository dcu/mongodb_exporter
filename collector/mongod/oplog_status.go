package collector_mongod

import (
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func GetCollectionSizeGB(db string, collection string, session *mgo.Session) (float64) {
    var collStats map[string]interface{}
    err := session.DB(db).Run(bson.D{{"collStats", collection }}, &collStats)
    if err != nil {
        glog.Error("Error getting collection stats!")
    }

    var size int
    size = collStats["size"].(int)
    return float64(size)/1024/1024/1024
}

func GetOplogSizeGB(session *mgo.Session) (float64) {
    return GetCollectionSizeGB("local", "oplog.rs", session)
}

func ParseBsonMongoTimestamp(timestamp bson.MongoTimestamp) (int32) {
    ts := (timestamp >> 32)
    return int32(ts)
}

func GetOplogLengthSecs(session *mgo.Session) (int32) {
    col := session.DB("local").C("oplog.rs")

    var head map[string]interface{}
    err := col.Find(bson.M{}).Sort("$natural").One(&head)
    if err != nil {
        glog.Error("Error getting head of oplog.rs!")
    } 

    var tail map[string]interface{}
    err = col.Find(bson.M{}).Sort("-$natural").One(&tail)
    if err != nil {
        glog.Error("Error getting tail of oplog.rs!")
    }

    head_ts := ParseBsonMongoTimestamp(head["ts"].(bson.MongoTimestamp))
    tail_ts := ParseBsonMongoTimestamp(tail["ts"].(bson.MongoTimestamp))
   
    return tail_ts - head_ts
}
