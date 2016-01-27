package collector_mongos

import (
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

    var result float64 = 0
    if balancerConfig["stopped"] != nil {
        balancerStopped := balancerConfig["stopped"].(bool)
        if balancerStopped == false {
            result = 1
        }
    }

    return result
}

func GetAllShardChunkInfo(session *mgo.Session) (map[string]int64) {
    var result []map[string]int64
    err := session.DB("config").C("chunks").Pipe([]bson.M{{ "$group" : bson.M{ "_id" : "$shard", "count" : bson.M{ "$sum" : 1  } } }}).All(&result)
    if err != nil {
        glog.Error("Could not find shard chunk info!")
    }

    shardChunkCounts := make(map[string]int64)
    for _, element := range result {
        shard := string(element["_id"])
        shardChunkCounts[shard] = int64(element["count"])
    }

    return shardChunkCounts
}

func IsClusterBalanced(session *mgo.Session) (float64) {
    // Different thresholds based on size
    // http://docs.mongodb.org/manual/core/sharding-internals/#sharding-migration-thresholds
    var threshold int64
    totalChunkCount := GetTotalChunks(session)
    if totalChunkCount < 20 {
        threshold = 2
    } else if totalChunkCount < 80 && totalChunkCount > 21 {
        threshold = 4
    } else {
        threshold = 8
    }

    var minChunkCount int64 = -1
    var maxChunkCount int64 = 0
    shardChunkInfoAll := GetAllShardChunkInfo(session)
    for _, chunkCount := range shardChunkInfoAll {
        if chunkCount > maxChunkCount {
            maxChunkCount = chunkCount
        }
        if minChunkCount == -1 || chunkCount < minChunkCount {
            minChunkCount = chunkCount
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
    IsBalanced		float64	
    BalancerEnabled	float64
}

func (status *BalancerStats) Export(ch chan<- prometheus.Metric) {
    balancerInfo.WithLabelValues("is_balanced").Set(status.IsBalanced)
    balancerInfo.WithLabelValues("balancer_on").Set(status.BalancerEnabled)
    balancerInfo.Collect(ch)
}

func GetBalancerStatus(session *mgo.Session) *BalancerStats {
    results := &BalancerStats{}

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)

    results.IsBalanced = IsClusterBalanced(session)
    results.BalancerEnabled = IsBalancerEnabled(session)

    return results
}
