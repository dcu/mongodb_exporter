package shared

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Groups contains all tracked groups
	Groups = make(map[string]*Group)
)

// Group is the definition of a group.
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

// FindOrCreateGroup finds or creates a group given a name
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

// CollectAllGroups go for each group collecting the metrics.
func CollectAllGroups(ch chan<- prometheus.Metric) {
	for _, group := range Groups {
		group.Collect(ch)
	}
}

// Export exports the given field and value to prometheus.
func (group *Group) Export(fieldName string, value float64) {
	groupDesc := GroupsDesc[group.DescName]

	if groupDesc.Type == "metrics" {
		group.trackField(fieldName, value)
	} else {
		group.trackFieldsVec(groupDesc, []string{fieldName}, value)
	}
}

// ExportWithLabels exports the field given a list of labels.
func (group *Group) ExportWithLabels(labels []string, value float64) {
	groupDesc := GroupsDesc[group.DescName]
	group.trackFieldsVec(groupDesc, labels, value)
}

// Collect collects the metrics for the given group.
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

// Describe describes the group
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

// GetGauge gets a gauge metric for the given name.
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

// GetCounter gets a counter metric for the given name.
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

// GetSummary gets a summary metric for the given name.
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

func (group *Group) trackFieldsVec(groupDesc *GroupDesc, labels []string, value float64) {
	glog.Infof("Setting %s(%s)=%f (%s=%v)", group.Name, groupDesc.Type, value, groupDesc.Labels, labels)
	switch groupDesc.Type {
	case "counter_vec":
		{
			vector := group.getCounterVec(group.Name, groupDesc.Help, groupDesc.Labels)
			vector.WithLabelValues(labels...).Set(value)
		}
	case "gauge_vec":
		{
			vector := group.getGaugeVec(group.Name, groupDesc.Help, groupDesc.Labels)
			vector.WithLabelValues(labels...).Set(value)
		}
	case "summary_vec":
		{
			vector := group.getSummaryVec(group.Name, groupDesc.Help, groupDesc.Labels)
			vector.WithLabelValues(labels...).Observe(value)
		}
	default:
		{
			panic("Unknown metadata type: " + groupDesc.Type)
		}
	}
}
