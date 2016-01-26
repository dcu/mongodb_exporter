package collector_mongos

import (
    "time"
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var (
    balancerInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
            Namespace: Namespace,
            Name:      "balancer_info",
            Help:      "TBD",
    }, []string{"type"})
)

func IsBalancerEnabled(session *mgo.Session) (float64) {
    var balancerConfig map[string]interface{}
    err := session.DB("config").C("settings").Find(bson.M{ "_id" : "balancer" }).Select(bson.M{ "_id" : 0 }).One(&balancerConfig)
    if err != nil {
        glog.Error("Could not find balancer settings in 'config.settings'!")
    }

    balancerStopped := balancerConfig["stopped"].(bool)
    if balancerStopped == false {
        return 1
    }
    return 0
}

func GetTotalShards(session *mgo.Session) (float64) {
    shardCount, err := session.DB("config").C("shards").Find(bson.M{}).Count()
    if err != nil {
        glog.Error("Could not find shard information in 'config.settings'!")
    }
    return float64(shardCount)
}

func GetTotalChunks(session *mgo.Session) (float64) {
    chunkCount, err := session.DB("config").C("chunks").Find(bson.M{}).Count()
    if err != nil {
        glog.Error("Could not find chunk information in 'config.chunks'!")
    }
    return float64(chunkCount)
}

func GetTotalShardedDatabases(session *mgo.Session) (float64) {
    dbCount, err := session.DB("config").C("databases").Find(bson.M{ "partitioned" : true }).Count()
    if err != nil {
        glog.Error("Could not find database information in 'config.databases'!")
    }
    return float64(dbCount)
}

func GetTotalShardedCollections(session *mgo.Session) (float64) {
    collCount, err := session.DB("config").C("collections").Find(bson.M{ "dropped" : false }).Count()
    if err != nil {
        glog.Error("Could not find collection information in 'config.collections'!")
    }
    return float64(collCount)
}

func TwentyFourHoursAgo() (time.Time) {
    return time.Now().Add(-24 * time.Hour)
}

func GetBalancerRoundCount24h(errorOccured bool, session *mgo.Session) (float64) {
    findQuery := bson.M{ "what" : "balancer.round", "details.errorOccured" : errorOccured, "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    roundCount, err := session.DB("config").C("actionlog").Find(findQuery).Count()
    if err != nil {
        glog.Error("Could not find balancer round info in 'config.actionlog'!")
    }
    return float64(roundCount)
}

func GetFailedBalancerRoundCount24h(session *mgo.Session) (float64) {
    roundCount := GetBalancerRoundCount24h(true, session)
    return float64(roundCount)
}

func GetSplitCount24h(session *mgo.Session) (float64) {
    findQuery := bson.M{ "what" : "split", "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    splitCount, err := session.DB("config").C("changelog").Find(findQuery).Count()
    if err != nil {
        glog.Error("Could not find split information in 'config.changelog'!")
    }
    return float64(splitCount)
}

func GetMoveChunkStartCount24h(session *mgo.Session) (float64) {
    findQuery := bson.M{ "what" : "moveChunk.start", "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    moveChunkStartCount, err := session.DB("config").C("changelog").Find(findQuery).Count()
    if err != nil {
        glog.Error("Could not find moveChunk.start info in 'config.changelog'!")
    }
    return float64(moveChunkStartCount)
}

func GetAllShardChunkInfo(session *mgo.Session) (map[string]float64) {
    var result []map[string]int64
    err := session.DB("config").C("chunks").Pipe([]bson.M{{ "$group" : bson.M{ "_id" : "$shard", "count" : bson.M{ "$sum" : 1  } } }}).All(&result)
    if err != nil {
        glog.Error("Could not find shard chunk info!")
    }

    shardChunkCounts := make(map[string]float64)
    for _, element := range result {
        shard := string(element["_id"])
        shardChunkCounts[shard] = float64(element["count"])
    }

    return shardChunkCounts
}

func IsClusterBalanced(session *mgo.Session) (float64) {
    // Different thresholds based on size
    // http://docs.mongodb.org/manual/core/sharding-internals/#sharding-migration-thresholds
    var threshold float64
    totalChunkCount := GetTotalChunks(session)
    if totalChunkCount < 20 {
        threshold = 2
    } else if totalChunkCount < 80 && totalChunkCount > 21 {
        threshold = 4
    } else {
        threshold = 8
    }

    var minChunkCount float64 = -1
    var maxChunkCount float64 = 0
    shardChunkInfoAll := GetAllShardChunkInfo(session)
    for _, chunkCount := range shardChunkInfoAll {
        if chunkCount > maxChunkCount {
            maxChunkCount = chunkCount
        } else {
            if minChunkCount < float64(0) {
                minChunkCount = chunkCount
            } else if chunkCount < minChunkCount {
                minChunkCount = chunkCount
            }
        }
    }

    // return true if the difference between the min and max is < the thresold
    chunkDifference := maxChunkCount - minChunkCount
    if chunkDifference < threshold {
        return 1
    }

    return 0
}

type BalancerStats struct {
    IsBalanced			float64	
    BalancerEnabled		float64
    TotalShards			float64
    TotalChunks			float64
    TotalDatabases		float64
    TotalCollections		float64
    BalancerRoundFailureLast24h	float64
    SplitCountLast24h		float64
    MoveChunkStartCountLast24h	float64
}

func (status *BalancerStats) Export(ch chan<- prometheus.Metric) {
    balancerInfo.WithLabelValues("is_balanced").Set(status.IsBalanced)
    balancerInfo.WithLabelValues("balancer_on").Set(status.BalancerEnabled)
    balancerInfo.WithLabelValues("total_shards").Set(status.TotalShards)
    balancerInfo.WithLabelValues("total_chunks").Set(status.TotalChunks)
    balancerInfo.WithLabelValues("total_dbs").Set(status.TotalDatabases)
    balancerInfo.WithLabelValues("total_cols").Set(status.TotalCollections)
    balancerInfo.WithLabelValues("failed_balances_24h").Set(status.BalancerRoundFailureLast24h)
    balancerInfo.WithLabelValues("splits_24h").Set(status.SplitCountLast24h)
    balancerInfo.WithLabelValues("move_chunk_starts_24h").Set(status.MoveChunkStartCountLast24h)
    balancerInfo.Collect(ch)
}

func GetBalancerStatus(session *mgo.Session) *BalancerStats {
    results := &BalancerStats{}

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)

    results.IsBalanced = IsClusterBalanced(session)
    results.BalancerEnabled = IsBalancerEnabled(session)
    results.TotalShards = GetTotalShards(session)
    results.TotalChunks = GetTotalChunks(session)
    results.TotalDatabases = GetTotalShardedDatabases(session)
    results.TotalCollections = GetTotalShardedCollections(session)
    results.BalancerRoundFailureLast24h = GetFailedBalancerRoundCount24h(session)
    results.SplitCountLast24h = GetSplitCount24h(session)
    results.MoveChunkStartCountLast24h = GetMoveChunkStartCount24h(session)

    return results
}
