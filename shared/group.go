package shared

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Groups = make(map[string]*Group)
)

type Group struct {
	Name        string
	DescName    string
	Counters    map[string]prometheus.Counter
	CounterVecs map[string]*prometheus.CounterVec
	Gauges      map[string]prometheus.Gauge
	GaugeVecs   map[string]*prometheus.GaugeVec
	Summaries   map[string]prometheus.Summary
	SummaryVecs map[string]*prometheus.SummaryVec
}

func FindOrCreateGroup(name string) *Group {
	name = SnakeCase(name)
	group := Groups[name]

	if group == nil {
		group = &Group{
			Name:        name,
			DescName:    name,
			Counters:    make(map[string]prometheus.Counter),
			CounterVecs: make(map[string]*prometheus.CounterVec),
			Gauges:      make(map[string]prometheus.Gauge),
			GaugeVecs:   make(map[string]*prometheus.GaugeVec),
			Summaries:   make(map[string]prometheus.Summary),
			SummaryVecs: make(map[string]*prometheus.SummaryVec),
		}
		Groups[name] = group
	}

	return group
}

func CollectAllGroups(ch chan<-prometheus.Metric) {
	for _, group := range Groups {
		group.Collect(ch)
	}
}

func (group *Group) Export(fieldName string, value float64) {
	fields := GroupFields(group.DescName)
	groupType := fields["metadata"].Type

	if groupType == "metrics" {
		group.trackField(fieldName, value)
	} else {
		group.trackFieldsVec(fields, fieldName, value)
	}
}

func (group *Group) Collect(ch chan<- prometheus.Metric) {
	if ch == nil {
		return
	}

	for _, counter := range group.Counters {
		counter.Collect(ch)
	}
	for _, counter := range group.CounterVecs {
		counter.Collect(ch)
	}

	for _, gauge := range group.Gauges {
		gauge.Collect(ch)
	}
	for _, gauge := range group.GaugeVecs {
		gauge.Collect(ch)
	}

	for _, summary := range group.Summaries {
		summary.Collect(ch)
	}
	for _, summary := range group.SummaryVecs {
		summary.Collect(ch)
	}
}

func (group *Group) Describe(ch chan<- *prometheus.Desc) {
	for _, counter := range group.Counters {
		counter.Describe(ch)
	}
	for _, counter := range group.CounterVecs {
		counter.Describe(ch)
	}

	for _, gauge := range group.Gauges {
		gauge.Describe(ch)
	}
	for _, gauge := range group.GaugeVecs {
		gauge.Describe(ch)
	}

	for _, summary := range group.Summaries {
		summary.Describe(ch)
	}
	for _, summary := range group.SummaryVecs {
		summary.Describe(ch)
	}
}

func (group *Group) GetGauge(name string, description string) prometheus.Gauge {
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

func (group *Group) getGaugeVec(name string, description string, labelNames []string) *prometheus.GaugeVec {
	gaugeVec := group.GaugeVecs[name]

	if gaugeVec == nil {
		gaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mongodb",
			Name:      name,
			Help:      description,
		}, labelNames)
		group.GaugeVecs[name] = gaugeVec
	}

	return gaugeVec
}

func (group *Group) GetCounter(name string, description string) prometheus.Counter {
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

func (group *Group) getCounterVec(name string, description string, labelNames []string) *prometheus.CounterVec {
	counter := group.CounterVecs[name]

	if counter == nil {
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "mongodb",
				Name:      name,
				Help:      description,
			},
			labelNames,
		)
		group.CounterVecs[name] = counter
	}

	return counter
}

func (group *Group) GetSummary(name string, description string) prometheus.Summary {
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

func (group *Group) getSummaryVec(name string, description string, labelNames []string) *prometheus.SummaryVec {
	summaryVec := group.SummaryVecs[name]

	if summaryVec == nil {
		summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: "mongodb",
			Name:      name,
			Help:      description,
		}, labelNames)
		group.SummaryVecs[name] = summaryVec
	}

	return summaryVec
}

func (group *Group) trackField(fieldName string, value float64) {
	fieldDesc := GroupField(group.DescName, fieldName)
	collectorType := fieldDesc.Type

	glog.Infof("Setting %s(metrics)=%f (%s,%s)", group.Name, value, collectorType, fieldName)
	switch collectorType {
	case "counter":
		{
			counter := group.GetCounter(fieldName, fieldDesc.Help)
			counter.Set(value)
		}
	case "gauge":
		{
			gauge := group.GetGauge(fieldName, fieldDesc.Help)
			gauge.Set(value)
		}
	case "summary":
		{
			summary := group.GetSummary(fieldName, fieldDesc.Help)
			summary.Observe(value)
		}
	}
}

func (group *Group) trackFieldsVec(fields GroupFieldsMap, fieldName string, value float64) {
	metadata := fields["metadata"]

	if fields[fieldName] == nil {
		panic("Label not declared in groups.yml file: "+group.Name + "."+fieldName)
	}

	glog.Infof("Setting %s(%s)=%f (%s=%s)", group.Name, metadata.Type, value, metadata.Label, fieldName)
	switch metadata.Type {
	case "counter_vec":
		{
			vector := group.getCounterVec(group.Name, metadata.Help, []string{metadata.Label})
			vector.WithLabelValues(fieldName).Set(value)
		}
	case "gauge_vec":
		{
			vector := group.getGaugeVec(group.Name, metadata.Help, []string{metadata.Label})
			vector.WithLabelValues(fieldName).Set(value)
		}
	case "summary_vec":
		{
			vector := group.getSummaryVec(group.Name, metadata.Help, []string{metadata.Label})
			vector.WithLabelValues(fieldName).Observe(value)
		}
	default:
		{
			panic("Unknown metadata type: " + metadata.Type)
		}
	}
}
