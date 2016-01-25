package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func GetCollectionSizeGB(db string, collection string, session *mgo.Session) (float64) {
    var collStats map[string]interface{}
    err := session.DB(db).Run(bson.D{{"collStats", collection }}, &collStats)
    if err != nil {
        panic(err)
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
        panic(err)
    } 

    var tail map[string]interface{}
    err = col.Find(bson.M{}).Sort("-$natural").One(&tail)
    if err != nil {
        panic(err)
    }

    head_ts := ParseBsonMongoTimestamp(head["ts"].(bson.MongoTimestamp))
    tail_ts := ParseBsonMongoTimestamp(tail["ts"].(bson.MongoTimestamp))
   
    return tail_ts - head_ts
}

func main() {
    uri := "mongodb://localhost:27017"
    session, err := mgo.Dial(uri)

    if err != nil {
        fmt.Println("Failed to get collection stats.")
    }

    fmt.Println("oplog size:", GetOplogSizeGB(session), "gb")
    fmt.Println("oplog length:", GetOplogLengthSecs(session), "sec(s)")
}
