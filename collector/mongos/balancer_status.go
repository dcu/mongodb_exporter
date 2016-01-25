package main

import (
    "time"
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func GetBalancerEnabled(session *mgo.Session) (bool) {
    var balancerConfig map[string]interface{}
    err := session.DB("config").C("settings").Find(bson.M{ "_id" : "balancer" }).Select(bson.M{ "_id" : 0 }).One(&balancerConfig)
    if err != nil {
        panic(err)
    }

    balancerStopped := balancerConfig["stopped"].(bool)
    if balancerStopped == false {
        return true
    }
    return false
}

func GetTotalShards(session *mgo.Session) (int) {
    shardCount, err := session.DB("config").C("shards").Find(bson.M{}).Count()
    if err != nil {
        panic(err)
    }
    return shardCount
}

func GetTotalChunks(session *mgo.Session) (int) {
    chunkCount, err := session.DB("config").C("chunks").Find(bson.M{}).Count()
    if err != nil {
        panic(err)
    }
    return chunkCount
}

func GetTotalShardedDatabases(session *mgo.Session) (int) {
    dbCount, err := session.DB("config").C("databases").Find(bson.M{ "partitioned" : true }).Count()
    if err != nil {
        panic(err)
    }
    return dbCount
}

func GetTotalShardedCollections(session *mgo.Session) (int) {
    collCount, err := session.DB("config").C("collections").Find(bson.M{ "dropped" : false }).Count()
    if err != nil {
        panic(err)
    }
    return collCount
}

func TwentyFourHoursAgo() (time.Time) {
    return time.Now().Add(-24 * time.Hour)
}

func GetBalancerRoundCount24h(errorOccured bool, session *mgo.Session) (int) {
    findQuery := bson.M{ "what" : "balancer.round", "details.errorOccured" : errorOccured, "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    roundCount, err := session.DB("config").C("actionlog").Find(findQuery).Count()
    if err != nil {
        panic(err)
    }
    return roundCount
}

func GetFailedBalancerRoundCount24h(session *mgo.Session) (int) {
    return GetBalancerRoundCount24h(true, session)
}

func GetOkBalancerRoundCount24h(session *mgo.Session) (int) {
    return GetBalancerRoundCount24h(false, session)
}

func GetSplitCount24h(session *mgo.Session) (int) {
    findQuery := bson.M{ "what" : "split", "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    splitCount, err := session.DB("config").C("changelog").Find(findQuery).Count()
    if err != nil {
        panic(err)
    }
    return splitCount
}

func GetMoveChunkStartCount24h(session *mgo.Session) (int) {
    findQuery := bson.M{ "what" : "moveChunk.start", "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    moveChunkStartCount, err := session.DB("config").C("changelog").Find(findQuery).Count()
    if err != nil {
        panic(err)
    }
    return moveChunkStartCount
}

type ShardChunkInfo struct {
    Count int
    PercentTot float64
}

func GetAllShardChunkInfo(session *mgo.Session) (map[string]ShardChunkInfo) {
    var result []map[string]interface{}
    err := session.DB("config").C("chunks").Pipe([]bson.M{{ "$group" : bson.M{ "_id" : "$shard", "count" : bson.M{ "$sum" : 1  } } }}).All(&result)
    if err != nil {
        panic(err)
    }

    totalChunks := GetTotalChunks(session)
    shardChunkCounts := make(map[string]ShardChunkInfo)
    for _, element := range result {
        shard := element["_id"].(string)
        count := element["count"].(int)
        ptot := float64(count) / float64(totalChunks) * 100
        shardChunkCounts[shard] = ShardChunkInfo{ Count: count, PercentTot: ptot }
    }

    return shardChunkCounts
}

func IsClusterBalanced(session *mgo.Session) (bool) {

    // Different thresholds based on size
    // http://docs.mongodb.org/manual/core/sharding-internals/#sharding-migration-thresholds
    var threshold int
    totalChunkCount := GetTotalChunks(session)
    if totalChunkCount < 20 {
        threshold = 2
    } else if totalChunkCount < 80 && totalChunkCount > 21 {
        threshold = 4
    } else {
        threshold = 8
    }

    var minChunkCount int = -1
    var maxChunkCount int = 0
    shardChunkInfoAll := GetAllShardChunkInfo(session)
    for _, shardChunkInfo := range shardChunkInfoAll {
        chunkCount := shardChunkInfo.Count
        if chunkCount > maxChunkCount {
            maxChunkCount = chunkCount
        } else {
            if minChunkCount < 0 {
                minChunkCount = chunkCount
            } else if chunkCount < minChunkCount {
                minChunkCount = chunkCount
            }
        }
    }

    // return true if the difference between the min and max is < the thresold
    chunkDifference := maxChunkCount - minChunkCount
    if chunkDifference < threshold {
        return true
    }

    return false
}

func main() {
    uri := "mongodb://localhost:27018"
    session, err := mgo.Dial(uri)
    if err != nil {
        fmt.Println("Failed to get collection stats.")
    }

    fmt.Println("balancer enabled:", GetBalancerEnabled(session))
    fmt.Println("cluster balanced:", IsClusterBalanced(session))
    fmt.Println("total shards:", GetTotalShards(session))
    fmt.Println("total chunks:", GetTotalChunks(session))
    fmt.Println("total sharding-enabled dbs:", GetTotalShardedDatabases(session))
    fmt.Println("total sharded collections:", GetTotalShardedCollections(session))
    fmt.Println("failed balancer rounds in last 24hrs:", GetFailedBalancerRoundCount24h(session))
    fmt.Println("moveChunk.start ops in last 24hrs:", GetMoveChunkStartCount24h(session))
    fmt.Println("chunk splits in last 24hrs:", GetSplitCount24h(session))

    session.Close()
}
