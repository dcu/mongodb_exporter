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

type BalancerChangelogAggregationId struct {
    Event	string	`bson:"event"`
    Note	string	`bson:"note"`
}

type BalancerChangelogAggregationResult struct {
    Id		*BalancerChangelogAggregationId	`bson:"_id"`
    Count 	float64				`bson:"count"`
}

type BalancerChangelogStats struct {
    MoveChunkStart              float64
    MoveChunkFromSuccess        float64
    MoveChunkFromFailed         float64
    MoveChunkToSuccess          float64
    MoveChunkToFailed           float64
    MoveChunkCommit             float64
    Split                       float64
    MultiSplit                  float64
    ShardCollection             float64
    ShardCollectionStart        float64
    AddShard                    float64
}

func GetBalancerChangelogStats24hr(session *mgo.Session, showErrors bool) *BalancerChangelogStats {
    var qresults []BalancerChangelogAggregationResult
    coll  := session.DB("config").C("changelog")
    match := bson.M{ "time" : bson.M{ "$gt" : time.Now().Add(-24 * time.Hour) } }
    group := bson.M{ "_id" : bson.M{ "event" : "$what", "note" : "$details.note" }, "count" : bson.M{ "$sum" : 1 } }

    err := coll.Pipe([]bson.M{ { "$match" : match }, { "$group" : group } }).All(&qresults)
    if err != nil {
        glog.Error("Error executing aggregation on 'config.changelog'!")
    }

    results := &BalancerChangelogStats{}
    for _, stat := range qresults {
        if stat.Id.Event == "moveChunk.start" {
            results.MoveChunkStart = stat.Count
        } else if stat.Id.Event == "moveChunk.to" {
            if stat.Id.Note == "success" {
                results.MoveChunkToSuccess = stat.Count
            } else {
                results.MoveChunkToFailed = stat.Count
            }
        } else if stat.Id.Event == "moveChunk.from" {
            if stat.Id.Note == "success" {
                results.MoveChunkFromSuccess = stat.Count
            } else {
                results.MoveChunkFromFailed = stat.Count
            }
        } else if stat.Id.Event == "moveChunk.commit" {
            results.MoveChunkCommit = stat.Count
        } else if stat.Id.Event == "addShard" {
            results.AddShard = stat.Count
        } else if stat.Id.Event == "shardCollection" {
            results.ShardCollection = stat.Count
        } else if stat.Id.Event == "shardCollection.start" {
            results.ShardCollectionStart = stat.Count
        } else if stat.Id.Event == "split" {
            results.Split = stat.Count
        } else if stat.Id.Event == "multi-split" {
            results.MultiSplit = stat.Count
        }
    }

    return results
}

func (status *BalancerChangelogStats) Export(ch chan<- prometheus.Metric) {
    balancerChangelogInfo.WithLabelValues("move_chunk_start").Set(status.MoveChunkStart)
    balancerChangelogInfo.WithLabelValues("move_chunk_to_success").Set(status.MoveChunkToSuccess)
    balancerChangelogInfo.WithLabelValues("move_chunk_to_failed").Set(status.MoveChunkToFailed)
    balancerChangelogInfo.WithLabelValues("move_chunk_from_success").Set(status.MoveChunkFromSuccess)
    balancerChangelogInfo.WithLabelValues("move_chunk_from_failed").Set(status.MoveChunkFromFailed)
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
    results := GetBalancerChangelogStats24hr(session, false)
    return results
}
