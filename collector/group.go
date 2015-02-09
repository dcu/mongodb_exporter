package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "reflect"
    "time"
)

type Group struct {
    Name      string
    Groups    map[string]*Group
    Counters  map[string]prometheus.Counter
    Gauges    map[string]prometheus.Gauge
    Summaries map[string]prometheus.Summary
}

func NewGroup(name string) (*Group) {
    group := &Group{
        Name: name,
        Counters: make(map[string]prometheus.Counter),
        Gauges: make(map[string]prometheus.Gauge),
        Summaries: make(map[string]prometheus.Summary),
        Groups: make(map[string]*Group),
    }
    return group
}

func (group *Group) Collect(object interface{}, fieldName string, ch chan<-prometheus.Metric) {
    value := reflect.Indirect(reflect.ValueOf(object))

    field, ok := value.Type().FieldByName(fieldName)
    fieldValue := value.FieldByName(fieldName)
    if !ok {
        panic("Field "+fieldName+" not found on group "+group.Name)
    }

    group.trackField(field, fieldValue, ch)
}

func (group *Group) Describe(ch chan<- *prometheus.Desc) {
    for _, counter := range group.Counters {
        counter.Describe(ch)
    }
    for _, gauge := range group.Gauges {
        gauge.Describe(ch)
    }
    for _, summary := range group.Summaries {
        summary.Describe(ch)
    }
}

func (group *Group) GetGauge(name string, description string) prometheus.Gauge {
    //println("Adding gauge", group.Name, name)

    gauge := group.Gauges[name]
    if gauge == nil {
        gauge = prometheus.NewGauge(prometheus.GaugeOpts{
            Namespace: "mongodb",
            Subsystem: group.Name,
            Name: name,
            Help: description,
        })
        group.Gauges[name] = gauge
    }

    return gauge
}

func (group *Group) GetCounter(name string, description string) prometheus.Counter {
    //println("Adding counter", group.Name, name)
    counter := group.Counters[name]

    if counter == nil {
        counter = prometheus.NewCounter(prometheus.CounterOpts{
            Namespace: "mongodb",
            Subsystem: group.Name,
            Name: name,
            Help: description,
        })
        group.Counters[name] = counter
    }

    return counter
}

func (group *Group) GetSummary(name string, description string) prometheus.Summary {
    //println("Adding summary", group.Name, name)
    summary := group.Summaries[name]

    if summary == nil {
        summary = prometheus.NewSummary(prometheus.SummaryOpts{
            Namespace: "mongodb",
            Subsystem: group.Name,
            Name: name,
            Help: description,
        })
        group.Summaries[name] = summary
    }

    return summary
}

func (group *Group) trackField(field reflect.StructField, fieldValue reflect.Value, ch chan<-prometheus.Metric) {
    valueToSet := getValueToSet(fieldValue)
    collectorType := field.Tag.Get("type")
    name := SnakeCase(field.Tag.Get("bson"))

    var collector prometheus.Collector
    switch(collectorType) {
        case "counter": {
            //println("Set", name, valueToSet)
            counter := group.GetCounter(name, "FIXME")
            counter.Set(valueToSet)
            collector = counter
        }
        case "gauge": {
            //println("Set", name, valueToSet)
            gauge := group.GetGauge(name, "FIXME")
            gauge.Set(valueToSet)
            collector = gauge
        }
        case "summary": {
            //println("Set", name, valueToSet)
            summary := group.GetGauge(name, "FIXME")
            summary.Set(valueToSet)
            collector = summary
        }
    }

    if ch != nil && collector != nil {
        collector.Collect(ch)
    }
}

func getValueToSet(fieldValue reflect.Value) float64 {
    var valueToSet float64
    if fieldValue.Kind() == reflect.Struct {
        valueToSet = getTimeFromFieldValue(fieldValue)
    } else {
        valueToSet = fieldValue.Float()
    }

    return valueToSet
}

func getTimeFromFieldValue(fieldValue reflect.Value) float64 {
    time, ok := fieldValue.Interface().(time.Time)

    var t float64
    if ok {
        t = float64(time.Unix())
    }

    return t
}

