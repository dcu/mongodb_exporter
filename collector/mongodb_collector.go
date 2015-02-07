package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "reflect"
)

type MongodbCollector struct {
    Groups []*Group
}

func NewMongodbCollector() *MongodbCollector {
    return &MongodbCollector{}
}

func (exporter *MongodbCollector) LoadGroups() {
    exporter.NewGroupFromStruct(reflect.TypeOf(DurTiming{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(DurStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(FlushStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(ConnectionStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(ExtraInfo{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(QueueStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(ClientStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(GlobalLockStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(IndexCounterStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(ReadWriteLockTimes{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(LockStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(OpcountStats{}))
    exporter.NewGroupFromStruct(reflect.TypeOf(MemStats{}))
}

func (exporter *MongodbCollector) NewGroup(name string) *Group {
    group := NewGroup(name)
    exporter.Groups = append(exporter.Groups, group)

    return group
}

func (exporter *MongodbCollector) NewGroupFromStruct(t reflect.Type) *Group {
    group := NewGroupFromStruct(t)
    exporter.Groups = append(exporter.Groups, group)

    return group
}

func (exporter *MongodbCollector) Describe(ch chan<- *prometheus.Desc) {
    for _,group := range exporter.Groups {
        group.Describe(ch)
    }
}

func (exporter *MongodbCollector) Collect(ch chan<- prometheus.Metric) {
    println("Collecting Server Status")
    //serverStatus := NewServerStatus()

    //connections.Set(serverStatus.Connections.Current)
    //connections.Collect(ch)
}


