package collector_mongos

import (
    "time"
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var (
    balancerChangelogInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
            Namespace: Namespace,
            Name:      "balancer_changelog",
            Help:      "Log event statistics for the MongoDB balancer",
    }, []string{"event"})
)

type BalancerChangelogAggregationResult struct {
    Event string	`bson:"_id"`
    Count float64	`bson:"count"`
}

type BalancerChangelogStats struct {
    MoveChunkStart              float64
    MoveChunkFrom               float64
    MoveChunkTo                 float64
    MoveChunkCommit             float64
    Split                       float64
    MultiSplit                  float64
    ShardCollection             float64
    ShardCollectionStart        float64
    AddShard                    float64
}

func GetBalancerChangelogStats24hr(session *mgo.Session) *BalancerChangelogStats {
    var qresults []BalancerChangelogAggregationResult
    coll  := session.DB("config").C("changelog")
    match := bson.M{ "time" : bson.M{ "$gt" : time.Now().Add(-24 * time.Hour) } }
    group := bson.M{ "_id" : "$what", "count" : bson.M{ "$sum" : 1 } }

    err := coll.Pipe([]bson.M{ { "$match" : match }, { "$group" : group } }).All(&qresults)
    if err != nil {
        glog.Error("Error executing aggregation on 'config.changelog'!")
    }

    results := &BalancerChangelogStats{}
    for _, stat := range qresults {
        if stat.Event == "moveChunk.start" {
            results.MoveChunkStart = stat.Count
        } else if stat.Event == "moveChunk.to" {
            results.MoveChunkTo = stat.Count
        } else if stat.Event == "moveChunk.from" {
            results.MoveChunkFrom = stat.Count
        } else if stat.Event == "moveChunk.commit" {
            results.MoveChunkCommit = stat.Count
        } else if stat.Event == "addShard" {
            results.AddShard = stat.Count
        } else if stat.Event == "shardCollection" {
            results.ShardCollection = stat.Count
        } else if stat.Event == "shardCollection.start" {
            results.ShardCollectionStart = stat.Count
        } else if stat.Event == "split" {
            results.Split = stat.Count
        } else if stat.Event == "multi-split" {
            results.MultiSplit = stat.Count
        }
    }

    return results
}

func (status *BalancerChangelogStats) Export(ch chan<- prometheus.Metric) {
    balancerChangelogInfo.WithLabelValues("move_chunk_start").Set(status.MoveChunkStart)
    balancerChangelogInfo.WithLabelValues("move_chunk_to").Set(status.MoveChunkTo)
    balancerChangelogInfo.WithLabelValues("move_chunk_from").Set(status.MoveChunkFrom)
    balancerChangelogInfo.WithLabelValues("move_chunk_commit").Set(status.MoveChunkCommit)
    balancerChangelogInfo.WithLabelValues("add_shard").Set(status.AddShard)
    balancerChangelogInfo.WithLabelValues("shard_collection").Set(status.ShardCollection)
    balancerChangelogInfo.WithLabelValues("shard_collection_start").Set(status.ShardCollectionStart)
    balancerChangelogInfo.WithLabelValues("split").Set(status.Split)
    balancerChangelogInfo.WithLabelValues("multi_split").Set(status.MultiSplit)
    balancerChangelogInfo.Collect(ch)
}

func GetBalancerChangelogStatus(session *mgo.Session) *BalancerChangelogStats {
    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)
    results := GetBalancerChangelogStats24hr(session)
    return results
}
