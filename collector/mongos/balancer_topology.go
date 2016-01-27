package collector_mongos

import (
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var (
    balancerTopoInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
            Namespace: Namespace,
            Name:      "balancer_topology",
            Help:      "Cluster topology statistics for the MongoDB balancer",
    }, []string{"type"})
)

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

type BalancerTopoStats struct {
    TotalShards		float64
    TotalChunks		float64
    TotalDatabases	float64
    TotalCollections	float64
}

func (status *BalancerTopoStats) Export(ch chan<- prometheus.Metric) {
    balancerTopoInfo.WithLabelValues("shards").Set(status.TotalShards)
    balancerTopoInfo.WithLabelValues("chunks").Set(status.TotalChunks)
    balancerTopoInfo.WithLabelValues("databases").Set(status.TotalDatabases)
    balancerTopoInfo.WithLabelValues("collections").Set(status.TotalCollections)
    balancerTopoInfo.Collect(ch)
}

func GetBalancerTopoStatus(session *mgo.Session) *BalancerTopoStats {
    results := &BalancerTopoStats{}

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)

    results.TotalShards = GetTotalShards(session)
    results.TotalChunks = GetTotalChunks(session)
    results.TotalDatabases = GetTotalShardedDatabases(session)
    results.TotalCollections = GetTotalShardedCollections(session)

    return results
}
