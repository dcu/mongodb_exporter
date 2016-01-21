package collector

import (
	//"github.com/dcu/mongodb_exporter/shared"
	"github.com/dcu/mongodb_exporter/collector/mongod"
	"github.com/dcu/mongodb_exporter/collector/mongos"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Namespace is the namespace of the metrics
	Namespace = "mongodb"
)

// MongodbCollectorOpts is the options of the mongodb collector.
type MongodbCollectorOpts struct {
	URI string
}

// MongodbCollector is in charge of collecting mongodb's metrics.
type MongodbCollector struct {
	Opts MongodbCollectorOpts
}

// NewMongodbCollector returns a new instance of a MongodbCollector.
func NewMongodbCollector(opts MongodbCollectorOpts) *MongodbCollector {
	exporter := &MongodbCollector{
		Opts: opts,
	}

	return exporter
}

// Describe describes all mongodb's metrics.
func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
	glog.Info("Describing groups")
	serverStatus := collector_mongod.GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Describe(ch)
	}
}

func connectMongo(uri string) *mgo.Session {
	session, err := mgo.Dial(uri)
	if err != nil {
		glog.Errorf("Cannot connect to server using url: %s", uri)
		return nil
	}

	session.SetMode(mgo.Eventual, true)
	session.SetSocketTimeout(0)
	defer func() {
		glog.Info("Closing connection to database.")
		session.Close()
	}()
	err = nil 
	return session, err 
}


// GetNodeType checks if the connected Session is a mongos, standalone, or replset,
// by looking at the result of calling isMaster.
func (session *mgo.Session) GetNodeType() (NodeType, error) {
	masterDoc := struct {
		SetName interface{} `bson:"setName"`
		Hosts   interface{} `bson:"hosts"`
		Msg     string      `bson:"msg"`
	}{}
	err = session.Run("isMaster", &masterDoc)
	if err != nil {
		return Unknown, err
	}

	if masterDoc.SetName != nil || masterDoc.Hosts != nil {
		return "replset", nil
	} else if masterDoc.Msg == "isdbgrid" {
		// isdbgrid is always the msg value when calling isMaster on a mongos
		// see http://docs.mongodb.org/manual/core/sharded-cluster-query-router/
		return "mongos", nil
	}
	return "mongod", nil
}

func AuthIfIneeded(session *mgo.Session, exporter *MongodbCollector){
	//Need to add some body here to try the user/pass/authdb opts if present, 
	//if not should run "ping" command to ensure access is allowed
}

// Collect collects all mongodb's metrics.
func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
    glog.Info("Collecting Server Status")
    session, err := connectMongo(exporter.Opts.URI)
    if err != nil:
		error()

    authIfNeeded(&session,exporter.Opts)
    nodeType := GetNodeType(session)

    switch nodeType {
        case 'mongos':
            collector_mongos.ServerStatus(ch, session)
        	collector_mongos.BalancingData(ch, session)
        case 'replset':
            collector_mongos.ServerStatus(ch, session)
           	collector_mongos.ElectionInfo(ch, session)
        	collector_mongos.OpLogInfo(ch, session)
        	collector_mongos.ReplicationInfo(ch, session)
        case "mongod":
            collector_mongod.ServerStatus(ch, session)
        case "arbiter":
        		continue
        default:
        	error()
    }
	//exporter.collectMongodServerStatus(ch)
	//exporter.collectMongosServerStatus(ch)
}

func (exporter *MongodbCollector) collectMongodServerStatus(ch chan<- prometheus.Metric) *collector_mongod.ServerStatus {
	serverStatus := collector_mongod.GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Export(ch)
	}

	return serverStatus
}

func (exporter *MongodbCollector) collectMongosServerStatus(ch chan<- prometheus.Metric) *collector_mongos.ServerStatus {
	serverStatus := collector_mongos.GetServerStatus(exporter.Opts.URI)

	if serverStatus != nil {
		serverStatus.Export(ch)
	}

	return serverStatus
}

