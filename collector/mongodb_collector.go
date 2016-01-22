package collector

import (
	//"github.com/dcu/mongodb_exporter/shared"
	"github.com/dcu/mongodb_exporter/collector/mongod"
	"github.com/dcu/mongodb_exporter/collector/mongos"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mgo.v2"
	"fmt"
	"os"
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
	session, err := connectMongo(exporter.Opts.URI)
    if err != nil{
		return
    }
	serverStatus := collector_mongos.GetServerStatus(session)

	if serverStatus != nil {
		serverStatus.Describe(ch)
	}
	defer session.Close()
}

func connectMongo(uri string)(*mgo.Session, error) {
	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
		glog.Errorf("Cannot connect to server using url: %s", uri)
		return nil,err
	}

	session.SetMode(mgo.Eventual, true)
	session.SetSocketTimeout(0)
	err = nil 
	return session,err
}

// GetNodeType checks if the connected Session is a mongos, standalone, or replset,
// by looking at the result of calling isMaster.
func GetNodeType(session *mgo.Session)(string, error) {
	masterDoc := struct {
		SetName interface{} `bson:"setName"`
		Hosts   interface{} `bson:"hosts"`
		Msg     string      `bson:"msg"`
	}{}
	err := session.Run("isMaster", &masterDoc)
	if err != nil {
		glog.Info("Got unknown node type\n")
		return "unknown", err
	}

	if masterDoc.SetName != nil || masterDoc.Hosts != nil {
		glog.Info("Got replset node type")
		return "replset", nil
	} else if masterDoc.Msg == "isdbgrid" {
		glog.Info("Got mongos node type\n")
		// isdbgrid is always the msg value when calling isMaster on a mongos
		// see http://docs.mongodb.org/manual/core/sharded-cluster-query-router/
		return "mongos", nil
	}
	glog.Info("defaulted to mongod node type\n")
	return "mongod", nil
}

// Collect collects all mongodb's metrics.
func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
    glog.Info("Collecting Server Status")
    session, err := connectMongo(exporter.Opts.URI)
    if err != nil{
		 glog.Error(fmt.Printf("We failed to connect to mongo with error of %s\n", err))
    }
    glog.Info("Passed connecting")
    nodeType,err := GetNodeType(session)
    if err != nil{
    	glog.Error(fmt.Printf("We run had a node type error of %s\n", err))
    }
	glog.Info(fmt.Printf("Passed nodeType with %s", nodeType))
    switch {
    	case nodeType == "mongos":
    		serverStatus := collector_mongos.GetServerStatus(session)
    		if serverStatus != nil {
				serverStatus.Export(ch)
			}
			//Need to build BalanceData
			//Need to build ChangeLogActions
		case nodeType == "mongod":
			serverStatus := collector_mongod.GetServerStatus(session)
    		if serverStatus != nil {
				serverStatus.Export(ch)
			}
		case nodeType == "replset":
			glog.Info("ReplicaSet stuff isnt setup yet!\n")
		default:
			glog.Info("No process for current node type no metrics printing!\n")
    }
    session.Close()
    /**
    switch nodeType {
    	case 'mongos':
            collector_mongos.ServerStatus(ch, session)
        	collector_mongos.BalancingData(ch, session)
        case 'replset':
            collector_mongos.ServerStatus(ch, session)
           	collector_mongos.ElectionInfo(ch, session)
        	collector_mongos.OpLogInfo(ch, session)
        	collector_mongos.ReplicationInfo(ch, session)
        case 'mongod':
            collector_mongod.ServerStatus(ch, session)
        case 'arbiter':
        		continue
        default:
        	error()
    }
    **/
	//exporter.collectMongodServerStatus(ch)
	//exporter.collectMongosServerStatus(ch)
}
/**
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
**/
