package collector_mongod

import (
    "time"
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var (
    replSetInfo = prometheus.NewCounterVec(prometheus.CounterOpts{
            Namespace: Namespace,
            Name:      "replset_info",
            Help:      "TBD",
    }, []string{"type"})
)

func GetReplSetConfig(session *mgo.Session) (map[string]interface{}) {
    var replSetConfig map[string]interface{}
    err := session.DB("admin").Run(bson.D{{"replSetGetConfig", 1}}, &replSetConfig)
    if err != nil {
        glog.Error("Error executing 'replSetGetConfig'!")
    }

    var result map[string]interface{}
    if replSetConfig["config"] != nil {
        result = replSetConfig["config"].(map[string]interface{})
    }

    return result
}

func GetReplSetMembers(session *mgo.Session) ([]interface{}) {
    replSetConfig := GetReplSetConfig(session)

    var result []interface{}
    if replSetConfig["members"] != nil {
        result = replSetConfig["members"].([]interface{})
    }

    return result
}

func GetReplSetMemberCount(session *mgo.Session) (float64) {
    replSetMembers := GetReplSetMembers(session)
    return float64(len(replSetMembers))
}

func GetReplSetMembersWithDataCount(session *mgo.Session) (float64) {
    replSetMembers := GetReplSetMembers(session)

    var membersWithDataCount int = 0
    if replSetMembers != nil {
        for _, member := range replSetMembers {
            memberInfo := member.(map[string]interface{})
            if memberInfo["arbiterOnly"] == false || memberInfo["health"] == 1 {
                membersWithDataCount = membersWithDataCount + 1
            }
        }
    }

    return float64(membersWithDataCount)
}

func GetReplSetMembersWithVotesCount(session *mgo.Session) (float64) {
    replSetMembers := GetReplSetMembers(session)

    var membersWithVotesCount int = 0
    if replSetMembers != nil {
        for _, member := range replSetMembers {
            memberInfo := member.(map[string]interface{})
            if memberInfo["votes"].(int) > 0 || memberInfo["health"] == 1 {
                membersWithVotesCount = membersWithVotesCount + 1
            }
        }
    }

    return float64(membersWithVotesCount)
}

func GetReplSetStatusInfo(session *mgo.Session) (map[string]interface{}) {
    var replSetStatus map[string]interface{}
    err := session.DB("admin").Run(bson.D{{"replSetGetStatus", 1}}, &replSetStatus)
    if err != nil {
        glog.Error("Error executing 'replSetStatus'!")
    }

    return replSetStatus
}

func GetReplSetStatusPrimary(session *mgo.Session) (map[string]interface{}) {
    replSetStatus := GetReplSetStatusInfo(session)

    var result map[string]interface{}
    if replSetStatus["members"] != nil {
        replSetStatusMembers := replSetStatus["members"].([]interface{})
        for _, member := range replSetStatusMembers {
            memberInfo := member.(map[string]interface{})
            if memberInfo["state"] == 1 {
                result = memberInfo
                break
            }
        }
    }

    return result
}

func GetReplStatusSelf(session *mgo.Session) (map[string]interface{}) {
    replSetStatus := GetReplSetStatusInfo(session)

    var result map[string]interface{}
    if replSetStatus["members"] != nil {
        replSetStatusMembers := replSetStatus["members"].([]interface{})
        for _, member := range replSetStatusMembers {
            memberInfo := member.(map[string]interface{})
            if memberInfo["self"] == true {
                result = memberInfo
                break
            }
        }
    }
    
    return result
}

func GetReplSetLagMs(session *mgo.Session) (float64) {
    memberInfo := GetReplStatusSelf(session)
    optimeNanoSelf := memberInfo["optimeDate"].(time.Time).UnixNano()

    // short-circuit the check if you're the Primary
    if memberInfo["state"] == 1 {
        return 0
    }

    var result float64 = -1
    replSetStatusPrimary := GetReplSetStatusPrimary(session)
    if replSetStatusPrimary["optimeDate"] != nil {
        optimeNanoPrimary := replSetStatusPrimary["optimeDate"].(time.Time).UnixNano()
        result = float64(optimeNanoPrimary - optimeNanoSelf)/1000000
    }

    return result
}

func GetReplSetLastElectionUnixTime(session *mgo.Session) (float64) {
    var result float64 = -1
    memberInfo := GetReplStatusSelf(session)
    if memberInfo["electionDate"] != nil {
        electionUnixTime := memberInfo["electionDate"].(time.Time).Unix()
        result = float64(electionUnixTime)
    } else {
        replSetPrimary := GetReplSetStatusPrimary(session)
        if replSetPrimary != nil {
            electionUnixTime := replSetPrimary["electionDate"].(time.Time).Unix()
            result = float64(electionUnixTime)
        }
    }

    return result
}

func GetReplSetMaxNode2NodePingMs(session *mgo.Session) (float64) {
    replSetStatus := GetReplSetStatusInfo(session)
    replSetStatusMembers := replSetStatus["members"].([]interface{})
    
    var maxNodePingMs float64 = -1
    for _, member := range replSetStatusMembers {
        memberInfo := member.(map[string]interface{})
        if memberInfo["pingMs"] != nil {
            pingMs := float64(memberInfo["pingMs"].(int))
            if pingMs > maxNodePingMs {
                maxNodePingMs = pingMs
            }
        }
    } 
    
    return maxNodePingMs
}

type ReplSetStats struct {
    Members		float64
    MembersWithData	float64
    MembersWithVotes	float64
    LagMs		float64
    MaxNode2NodePingMs	float64
    LastElection	float64
}

func(status *ReplSetStats) Export(ch chan<- prometheus.Metric) {
    replSetInfo.WithLabelValues("members").Set(status.Members)
    replSetInfo.WithLabelValues("members_w_data").Set(status.MembersWithData)
    replSetInfo.WithLabelValues("members_w_votes").Set(status.MembersWithVotes)
    replSetInfo.WithLabelValues("lag_ms").Set(status.LagMs)
    replSetInfo.WithLabelValues("max_n2n_ping_ms").Set(status.MaxNode2NodePingMs)
    replSetInfo.WithLabelValues("last_election").Set(status.LastElection)
    replSetInfo.Collect(ch)
}

func GetReplSetStatus(session *mgo.Session) *ReplSetStats {
  results := &ReplSetStats{}

  results.Members = GetReplSetMemberCount(session)
  results.MembersWithData = GetReplSetMembersWithDataCount(session)
  results.MembersWithVotes = GetReplSetMembersWithVotesCount(session)
  results.LagMs = GetReplSetLagMs(session)
  results.MaxNode2NodePingMs = GetReplSetMaxNode2NodePingMs(session)
  results.LastElection = GetReplSetLastElectionUnixTime(session)

  return results
}
