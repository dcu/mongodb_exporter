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
            Help:      "TBD",
    }, []string{"type"})
)

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

type BalancerChangelogStats struct {
    BalancerRoundFailure24h	float64
    SplitCount24h		float64
    MoveChunkStartCount24h	float64
}

func (status *BalancerChangelogStats) Export(ch chan<- prometheus.Metric) {
    balancerChangelogInfo.WithLabelValues("splits_24h").Set(status.SplitCount24h)
    balancerChangelogInfo.WithLabelValues("move_chunk_starts_24h").Set(status.MoveChunkStartCount24h)
    balancerChangelogInfo.WithLabelValues("move_chunk_failed_24h").Set(status.BalancerRoundFailure24h)
    balancerChangelogInfo.Collect(ch)
}

func GetBalancerChangelogStatus(session *mgo.Session) *BalancerChangelogStats {
    results := &BalancerChangelogStats{}

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)

    results.SplitCount24h = GetSplitCount24h(session)
    results.MoveChunkStartCount24h = GetMoveChunkStartCount24h(session)
    results.BalancerRoundFailure24h = GetFailedBalancerRoundCount24h(session)

    return results
}
