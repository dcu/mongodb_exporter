package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "reflect"
)

type Group struct {
    Name      string
    Counters  []prometheus.Counter
    Gauges    []prometheus.Gauge
    Summaries []prometheus.Summary
}

func NewGroup(name string) (*Group) {
    group := &Group{Name: name}
    return group
}

func NewGroupFromStruct(t reflect.Type) (*Group) {
	name := t.Name()
  group := NewGroup(name)

	for i := 0; i < t.NumField(); i++  {
		field := t.Field(i)
    group.AddMetricFromField(&field)
	}

  return group
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

func (group *Group) AddMetricFromField(field *reflect.StructField) {
		name := field.Tag.Get("bson")
		stat_type := field.Tag.Get("type")

    switch(stat_type) {
        case "gauge": {
            group.AddGauge(name, "FIXME")
        }
        case "counter": {
            group.AddCounter(name, "FIXME")
        }
        case "summary": {
            group.AddSummary(name, "FIXME")
        }
        default: {
            println("Warning:", name, "doesn't have the `type` tag")
        }
    }
}

func (group *Group) AddGauge(name string, description string) (*prometheus.Gauge) {
    gauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "mongodb",
        Subsystem: group.Name,
        Name: name,
        Help: description,
    })
    group.Gauges = append(group.Gauges, gauge)

    return &gauge
}

func (group *Group) AddCounter(name string, description string) (*prometheus.Counter) {
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "mongodb",
        Subsystem: group.Name,
        Name: name,
        Help: description,
    })
    group.Counters = append(group.Counters, counter)

    return &counter
}

func (group *Group) AddSummary(name string, description string) (*prometheus.Summary) {
    summary := prometheus.NewSummary(prometheus.SummaryOpts{
        Namespace: "mongodb",
        Subsystem: group.Name,
        Name: name,
        Help: description,
    })
    group.Summaries = append(group.Summaries, summary)

    return &summary
}

