package collector

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	subsystem = "replset"
	myState   = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "my_state",
		Help:      "An integer between 0 and 10 that represents the replica state of the current member",
	}, []string{"set"})

	myReplicaLag = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "my_replica_lag",
		Help:      "An integer shows the replication lag in seconds, -1 if no master found",
	}, []string{"set"})

	masterCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "master_count",
		Help:      "The number of master, any value except 1 means something wrong",
	}, []string{})

	term = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "term",
		Help:      "The election count for the replica set, as known to this replica set member",
	}, []string{"set"})
	numberOfMembers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "number_of_members",
		Help:      "The number of replica set mebers",
	}, []string{"set"})
	heartbeatIntervalMillis = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "heartbeat_interval_millis",
		Help:      "The frequency in milliseconds of the heartbeats",
	}, []string{"set"})
	memberHealth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_health",
		Help:      "This field conveys if the member is up (1) or down (0).",
	}, []string{"set", "name"})
	memberState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_state",
		Help:      "The value of state is an integer between 0 and 10 that represents the replica state of the member.",
	}, []string{"set", "name"})
	memberUptime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_uptime",
		Help:      "The uptime field holds a value that reflects the number of seconds that this member has been online.",
	}, []string{"set", "name"})
	memberOptimeDate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_optime_date",
		Help:      "The last entry from the oplog that this member applied.",
	}, []string{"set", "name"})
	memberElectionDate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_election_date",
		Help:      "The timestamp the node was elected as replica leader",
	}, []string{"set", "name"})
	memberLastHeartbeat = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_last_heartbeat",
		Help:      "The lastHeartbeat value provides an ISODate formatted date and time of the transmission time of last heartbeat received from this member",
	}, []string{"set", "name"})
	memberLastHeartbeatRecv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_last_heartbeat_recv",
		Help:      "The lastHeartbeatRecv value provides an ISODate formatted date and time that the last heartbeat was received from this member",
	}, []string{"set", "name"})
	memberPingMs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_ping_ms",
		Help:      "The pingMs represents the number of milliseconds (ms) that a round-trip packet takes to travel between the remote member and the local instance.",
	}, []string{"set", "name"})
	memberConfigVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_config_version",
		Help:      "The configVersion value is the replica set configuration version.",
	}, []string{"set", "name"})
	memberOptime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: subsystem,
		Name:      "member_optime",
		Help:      "Information regarding the last operation from the operation log that this member has applied.",
	}, []string{"set", "name"})
)

// ReplSetStatus keeps the data returned by the GetReplSetStatus method
type ReplSetStatus struct {
	Set                     string    `bson:"set"`
	Date                    time.Time `bson:"date"`
	MyState                 int32     `bson:"myState"`
	Term                    *int32    `bson:"term,omitempty"`
	HeartbeatIntervalMillis *float64  `bson:"heartbeatIntervalMillis,omitempty"`
	Members                 []Member  `bson:"members"`
}

// Member represents an array element of ReplSetStatus.Members
type Member struct {
	Name                 string      `bson:"name"`
	Self                 *bool       `bson:"self,omitempty"`
	Health               *int32      `bson:"health,omitempty"`
	State                int32       `bson:"state"`
	StateStr             string      `bson:"stateStr"`
	Uptime               float64     `bson:"uptime"`
	Optime               interface{} `bson:"optime"`
	OptimeDate           time.Time   `bson:"optimeDate"`
	ElectionTime         *time.Time  `bson:"electionTime,omitempty"`
	ElectionDate         *time.Time  `bson:"electionDate,omitempty"`
	LastHeartbeat        *time.Time  `bson:"lastHeartbeat,omitempty"`
	LastHeartbeatRecv    *time.Time  `bson:"lastHeartbeatRecv,omitempty"`
	LastHeartbeatMessage *string     `bson:"lastHeartbeatMessage,omitempty"`
	PingMs               *float64    `bson:"pingMs,omitempty"`
	SyncingTo            *string     `bson:"syncingTo,omitempty"`
	ConfigVersion        *int32      `bson:"configVersion,omitempty"`
}

// Export exports the replSetGetStatus stati to be consumed by prometheus
func (replStatus *ReplSetStatus) Export(ch chan<- prometheus.Metric) {
	myState.Reset()
	myReplicaLag.Reset()
	term.Reset()
	numberOfMembers.Reset()
	heartbeatIntervalMillis.Reset()
	memberState.Reset()
	memberHealth.Reset()
	memberUptime.Reset()
	memberOptimeDate.Reset()
	memberElectionDate.Reset()
	memberLastHeartbeat.Reset()
	memberLastHeartbeatRecv.Reset()
	memberPingMs.Reset()
	memberConfigVersion.Reset()
	masterCount.Reset()

	myState.WithLabelValues(replStatus.Set).Set(float64(replStatus.MyState))

	// new in version 3.2
	if replStatus.Term != nil {
		term.WithLabelValues(replStatus.Set).Set(float64(*replStatus.Term))
	}
	numberOfMembers.WithLabelValues(replStatus.Set).Set(float64(len(replStatus.Members)))

	// new in version 3.2
	if replStatus.HeartbeatIntervalMillis != nil {
		heartbeatIntervalMillis.WithLabelValues(replStatus.Set).Set(*replStatus.HeartbeatIntervalMillis)
	}

	var (
		primaryOpTime time.Time
		myOpTime      time.Time
	)
	mCount := 0
	for _, member := range replStatus.Members {
		ls := prometheus.Labels{
			"set":  replStatus.Set,
			"name": member.Name,
		}
		if member.State == 1 {
			primaryOpTime = member.OptimeDate
		}
		if member.Self != nil && *member.Self {
			myOpTime = member.OptimeDate
		}
		memberState.With(ls).Set(float64(member.State))
		if member.State == 1 {
			mCount += 1
		}

		// ReplSetStatus.Member.Health is not available on the node you're connected to
		if member.Health != nil {
			memberHealth.With(ls).Set(float64(*member.Health))
		}

		memberUptime.With(ls).Set(member.Uptime)

		memberOptimeDate.With(ls).Set(float64(member.OptimeDate.Unix()))

		// ReplSetGetStatus.Member.ElectionTime is only available on the PRIMARY
		if member.ElectionDate != nil {
			memberElectionDate.With(ls).Set(float64((*member.ElectionDate).Unix()))
		}
		if member.LastHeartbeat != nil {
			memberLastHeartbeat.With(ls).Set(float64((*member.LastHeartbeat).Unix()))
		}
		if member.LastHeartbeatRecv != nil {
			memberLastHeartbeatRecv.With(ls).Set(float64((*member.LastHeartbeatRecv).Unix()))
		}
		if member.PingMs != nil {
			memberPingMs.With(ls).Set(*member.PingMs)
		}
		if member.ConfigVersion != nil {
			memberConfigVersion.With(ls).Set(float64(*member.ConfigVersion))
		}
	}
	if !primaryOpTime.IsZero() && !myOpTime.IsZero() {
		myReplicaLag.WithLabelValues(replStatus.Set).Set(float64(primaryOpTime.Unix() - myOpTime.Unix()))
	} else {
		myReplicaLag.WithLabelValues(replStatus.Set).Set(-1.0)
	}
	masterCount.WithLabelValues().Set(float64(mCount))
	// collect metrics
	myState.Collect(ch)
	myReplicaLag.Collect(ch)
	term.Collect(ch)
	numberOfMembers.Collect(ch)
	heartbeatIntervalMillis.Collect(ch)
	memberState.Collect(ch)
	masterCount.Collect(ch)
	memberHealth.Collect(ch)
	memberUptime.Collect(ch)
	memberOptimeDate.Collect(ch)
	memberElectionDate.Collect(ch)
	memberLastHeartbeat.Collect(ch)
	memberLastHeartbeatRecv.Collect(ch)
	memberPingMs.Collect(ch)
	memberConfigVersion.Collect(ch)
}

// Describe describes the replSetGetStatus metrics for prometheus
func (replStatus *ReplSetStatus) Describe(ch chan<- *prometheus.Desc) {
	myState.Describe(ch)
	myReplicaLag.Describe(ch)
	term.Describe(ch)
	numberOfMembers.Describe(ch)
	heartbeatIntervalMillis.Describe(ch)
	memberState.Describe(ch)
	memberHealth.Describe(ch)
	memberUptime.Describe(ch)
	memberOptimeDate.Describe(ch)
	memberElectionDate.Describe(ch)
	memberLastHeartbeatRecv.Describe(ch)
	memberPingMs.Describe(ch)
	memberConfigVersion.Describe(ch)
}

// GetReplSetStatus returns the replica status info
func GetReplSetStatus(session *mgo.Session) *ReplSetStatus {
	result := &ReplSetStatus{}
	err := session.DB("admin").Run(bson.D{{"replSetGetStatus", 1}}, result)
	if err != nil {
		glog.Error("Failed to get replSet status.")
		return nil
	}
	return result
}
