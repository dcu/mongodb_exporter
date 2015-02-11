package shared

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Groups = make(map[string]*Group)
)

type Group struct {
	Name      string
	DescName  string
	Counters  map[string]prometheus.Counter
	Gauges    map[string]prometheus.Gauge
	Summaries map[string]prometheus.Summary
}

func FindOrCreateGroup(name string) *Group {
	name = SnakeCase(name)
	group := Groups[name]

	if group == nil {
		group = &Group{
			Name:      name,
			DescName:  name,
			Counters:  make(map[string]prometheus.Counter),
			Gauges:    make(map[string]prometheus.Gauge),
			Summaries: make(map[string]prometheus.Summary),
		}
		Groups[name] = group
	}

	return group
}

func (group *Group) Collect(fieldName string, value float64, ch chan<- prometheus.Metric) {
	fieldDesc := GroupField(group.DescName, fieldName)
	group.trackField(fieldName, fieldDesc, value, ch)
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
			Name:      name,
			Help:      description,
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
			Name:      name,
			Help:      description,
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
			Name:      name,
			Help:      description,
		})
		group.Summaries[name] = summary
	}

	return summary
}

func (group *Group) trackField(fieldName string, fieldDesc *FieldDesc, value float64, ch chan<- prometheus.Metric) {
	collectorType := fieldDesc.Type

	var collector prometheus.Collector
	switch collectorType {
	case "counter":
		{
			//println("Set", name, valueToSet)
			counter := group.GetCounter(fieldName, fieldDesc.Help)
			counter.Set(value)
			collector = counter
		}
	case "gauge":
		{
			//println("Set", name, valueToSet)
			gauge := group.GetGauge(fieldName, fieldDesc.Help)
			gauge.Set(value)
			collector = gauge
		}
	case "summary":
		{
			//println("Set", name, valueToSet)
			summary := group.GetGauge(fieldName, fieldDesc.Help)
			summary.Set(value)
			collector = summary
		}
	}

	if ch != nil && collector != nil {
		collector.Collect(ch)
	}
}
