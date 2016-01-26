package collector_mongod

import (
    "time"
    "github.com/golang/glog"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func GetReplSetConfig(session *mgo.Session) (map[string]interface{}) {
    var replSetConfig map[string]interface{}
    err := session.DB("admin").Run(bson.D{{"replSetGetConfig", 1}}, &replSetConfig)
    if err != nil {
        glog.Error("Error executing 'replSetGetConfig'!")
    }

    return replSetConfig["config"].(map[string]interface{})
}

func GetReplSetMembers(session *mgo.Session) ([]interface{}) {
    replSetConfig := GetReplSetConfig(session)
    return replSetConfig["members"].([]interface{})
}

func GetReplSetMemberCount(session *mgo.Session) (int) {
    replSetMembers := GetReplSetMembers(session)
    return len(replSetMembers)

}

func GetReplSetMembersWithDataCount(session *mgo.Session) (int) {
    replSetMembers := GetReplSetMembers(session)

    var membersWithDataCount int = 0
    for _, member := range replSetMembers {
        memberInfo := member.(map[string]interface{})
        if memberInfo["arbiterOnly"] == false {
            membersWithDataCount = membersWithDataCount + 1
        }
    } 

    return membersWithDataCount
}

func GetReplSetMembersWithVotesCount(session *mgo.Session) (int) {
    replSetMembers := GetReplSetMembers(session)

    var membersWithVotesCount int = 0
    for _, member := range replSetMembers {
        memberInfo := member.(map[string]interface{})
        if memberInfo["votes"].(int) > 0 {
            membersWithVotesCount = membersWithVotesCount + 1
        }
    }

    return membersWithVotesCount
}

func GetReplSetStatus(session *mgo.Session) (map[string]interface{}) {
    var replSetStatus map[string]interface{}
    err := session.DB("admin").Run(bson.D{{"replSetGetStatus", 1}}, &replSetStatus)
    if err != nil {
        glog.Error("Error executing 'replSetStatus'!")
    }

    return replSetStatus
}

func GetReplSetStatusPrimary(session *mgo.Session) (map[string]interface{}) {
    replSetStatus := GetReplSetStatus(session)
    replSetStatusMembers := replSetStatus["members"].([]interface{})

    for _, member := range replSetStatusMembers {
      memberInfo := member.(map[string]interface{})
      if memberInfo["state"] == 1 {
          return memberInfo
      }
    }

    glog.Error("Found no replSet member in Primary state!")
}

func GetReplStatusSelf(session *mgo.Session) (map[string]interface{}) {
    replSetStatus := GetReplSetStatus(session)
    replSetStatusMembers := replSetStatus["members"].([]interface{})

    for _, member := range replSetStatusMembers {
        memberInfo := member.(map[string]interface{})
        if memberInfo["self"] == true {
            return memberInfo
        }
    }
    
    glog.Error("Could not find myself in the replset config!")
}

func GetReplSetLagMs(session *mgo.Session) (float64) {
    memberInfo := GetReplStatusSelf(session)
    optimeNanoSelf := memberInfo["optimeDate"].(time.Time).UnixNano()

    // short-circuit the check if you're the Primary
    if memberInfo["state"] == 1 {
        return 0
    }

    replSetStatusPrimary := GetReplSetStatusPrimary(session)
    optimeNanoPrimary := replSetStatusPrimary["optimeDate"].(time.Time).UnixNano()

    return float64(optimeNanoPrimary - optimeNanoSelf)/1000000
}

func GetReplSetLastElectionUnixTime(session *mgo.Session) (int64) {
    memberInfo := GetReplStatusSelf(session)
    return memberInfo["electionDate"].(time.Time).Unix()
}

func GetReplSetMaxNode2NodePingMs(session *mgo.Session) (int) {
    replSetStatus := GetReplSetStatus(session)
    replSetStatusMembers := replSetStatus["members"].([]interface{})
    
    var maxNodePingMs int = 0
    for _, member := range replSetStatusMembers {
        memberInfo := member.(map[string]interface{})
        if memberInfo["pingMs"] != nil {
            pingMs := memberInfo["pingMs"].(int)
            if pingMs > maxNodePingMs {
                maxNodePingMs = pingMs
            }
        }
    } 
    
    return maxNodePingMs
}
