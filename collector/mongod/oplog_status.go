package collector_mongod

import (
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var (
    oplogInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
            Namespace: Namespace,
            Name:      "oplog_info",
            Help:      "TBD",
    }, []string{"type"})
)

func GetCollectionSizeGB(db string, collection string, session *mgo.Session) (float64) {
    var collStats map[string]interface{}
    err := session.DB(db).Run(bson.D{{"collStats", collection }}, &collStats)
    if err != nil {
        glog.Error("Error getting collection stats!")
    }

    var result float64 = -1
    if collStats["size"] != nil {
        size := collStats["size"].(int)
        result = float64(size)/1024/1024/1024
    }

    return result
}

func GetOplogSizeGB(session *mgo.Session) (float64) {
    return GetCollectionSizeGB("local", "oplog.rs", session)
}

func ParseBsonMongoTsToUnixTime(timestamp bson.MongoTimestamp) (int32) {
    return int32(timestamp >> 32)
}

func GetOplogLengthSecs(session *mgo.Session) (float64) {
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

    var result float64 = -1
    if head["ts"] != nil && tail["ts"] != nil {
        head_ts := ParseBsonMongoTsToUnixTime(head["ts"].(bson.MongoTimestamp))
        tail_ts := ParseBsonMongoTsToUnixTime(tail["ts"].(bson.MongoTimestamp))
        result = float64(tail_ts - head_ts)
    }

    return result
}

type OplogStats struct {
    LengthSec	float64
    SizeGB	float64
}

func (status *OplogStats) Export(ch chan<- prometheus.Metric) {
    oplogInfo.WithLabelValues("length_sec").Set(status.LengthSec)
    oplogInfo.WithLabelValues("size_gb").Set(status.SizeGB)
    oplogInfo.Collect(ch)
}

func GetOplogStatus(session *mgo.Session) *OplogStats {
    results := &OplogStats{}

    results.LengthSec = GetOplogLengthSecs(session)
    results.SizeGB = GetOplogSizeGB(session)

    return results
}
