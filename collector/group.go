package collector

import(
    "github.com/prometheus/client_golang/prometheus"
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


