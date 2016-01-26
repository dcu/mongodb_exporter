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

func IsBalancerEnabled(session *mgo.Session) (bool) {
    var balancerConfig map[string]interface{}
    err := session.DB("config").C("settings").Find(bson.M{ "_id" : "balancer" }).Select(bson.M{ "_id" : 0 }).One(&balancerConfig)
    if err != nil {
        glog.Error("Could not find balancer settings in 'config.settings'!")
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
        glog.Error("Could not find shard information in 'config.settings'!")
    }
    return shardCount
}

func GetTotalChunks(session *mgo.Session) (int) {
    chunkCount, err := session.DB("config").C("chunks").Find(bson.M{}).Count()
    if err != nil {
        glog.Error("Could not find chunk information in 'config.chunks'!")
    }
    return chunkCount
}

func GetTotalShardedDatabases(session *mgo.Session) (int) {
    dbCount, err := session.DB("config").C("databases").Find(bson.M{ "partitioned" : true }).Count()
    if err != nil {
        glog.Error("Could not find database information in 'config.databases'!")
    }
    return dbCount
}

func GetTotalShardedCollections(session *mgo.Session) (int) {
    collCount, err := session.DB("config").C("collections").Find(bson.M{ "dropped" : false }).Count()
    if err != nil {
        glog.Error("Could not find collection information in 'config.collections'!")
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
        glog.Error("Could not find balancer round info in 'config.actionlog'!")
    }
    return roundCount
}

func GetFailedBalancerRoundCount24h(session *mgo.Session) (int) {
    return GetBalancerRoundCount24h(true, session)
}

func GetSplitCount24h(session *mgo.Session) (int) {
    findQuery := bson.M{ "what" : "split", "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    splitCount, err := session.DB("config").C("changelog").Find(findQuery).Count()
    if err != nil {
        glog.Error("Could not find split information in 'config.changelog'!")
    }
    return splitCount
}

func GetMoveChunkStartCount24h(session *mgo.Session) (int) {
    findQuery := bson.M{ "what" : "moveChunk.start", "time" : bson.M{ "$gt" : TwentyFourHoursAgo() } }
    moveChunkStartCount, err := session.DB("config").C("changelog").Find(findQuery).Count()
    if err != nil {
        glog.Error("Could not find moveChunk.start info in 'config.changelog'!")
    }
    return moveChunkStartCount
}

func GetAllShardChunkInfo(session *mgo.Session) (map[string]int) {
    err := session.DB("config").C("chunks").Pipe([]bson.M{{ "$group" : bson.M{ "_id" : "$shard", "count" : bson.M{ "$sum" : 1  } } }}).All(&result)
    if err != nil {
        glog.Error("Could not find shard chunk info!")
    }

    shardChunkCounts := make(map[string]int)
    for _, element := range result {
        shard := element["_id"].(string)
        shardChunkCounts[shard] = element["count"].(int)
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
    for _, chunkCount := range shardChunkInfoAll {
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

type BalancerStatus struct {
    IsBalanced			bool	
    BalancerEnabled		bool
    TotalShards			int
    TotalChunks			int
    TotalDatabases		int
    TotalCollections		int
    BalancerRoundFailureLast24h	int
    SplitCountLast24h		int
    MoveChunkStartCountLast24h	int
}

func (status *BalancerStatus) Export(ch chan<- prometheus.Metric) {
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

func GetBalancerStatus(session *.mgo.Session) *BalancerStatus {
    results := &BalancerStatus{}

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
