package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics
type DocumentStats struct {
	Deleted  float64 `bson:"deleted"`
	Inserted float64 `bson:"inserted"`
	Returned float64 `bson:"returned"`
	Updated  float64 `bson:"updated"`
}

func (documentStats *DocumentStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("deleted", documentStats.Deleted, ch)
	group.Collect("inserted", documentStats.Inserted, ch)
	group.Collect("returned", documentStats.Returned, ch)
	group.Collect("updated", documentStats.Updated, ch)
}

type BenchmarkStats struct {
	Num         float64 `bson:"num"`
	TotalMillis float64 `bson:"totalMillis"`
}

func (benchmarkStats *BenchmarkStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("num", benchmarkStats.Num, ch)
	group.Collect("total_millis", benchmarkStats.TotalMillis, ch)
}

type GetLastErrorStats struct {
	Wtimeouts float64         `bson:"wtimeouts"`
	Wtime     *BenchmarkStats `bson:"wtime"`
}

func (getLastErrorStats *GetLastErrorStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("wtimeouts", getLastErrorStats.Wtimeouts, ch)
	getLastErrorStats.Wtime.Collect(groupName+"_wtime", ch)
}

type OperationStats struct {
	Fastmod      float64 `bson:"fastmod"`
	Idhack       float64 `bson:"idhack"`
	ScanAndOrder float64 `bson:"scanAndOrder"`
}

func (operationStats *OperationStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("fastmod", operationStats.Fastmod, ch)
	group.Collect("idhack", operationStats.Idhack, ch)
	group.Collect("scan_and_order", operationStats.ScanAndOrder, ch)
}

type QueryExecutorStats struct {
	Scanned        float64 `bson:"scanned"`
	ScannedObjects float64 `bson:"scannedObjects"`
}

func (queryExecutorStats *QueryExecutorStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("scanned", queryExecutorStats.Scanned, ch)
	group.Collect("scanned_objects", queryExecutorStats.ScannedObjects, ch)
}

type RecordStats struct {
	Moves float64 `bson:"moves"`
}

func (recordStats *RecordStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("moves", recordStats.Moves, ch)
}

type ApplyStats struct {
	Batches *BenchmarkStats `bson:"batches"`
	Ops     float64         `bson:"ops"`
}

func (applyStats *ApplyStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("ops", applyStats.Ops, ch)

	applyStats.Batches.Collect(groupName+"_batches", ch)
}

type BufferStats struct {
	Count        float64 `bson:"count"`
	MaxSizeBytes float64 `bson:"maxSizeBytes"`
	SizeBytes    float64 `bson:"sizeBytes"`
}

func (bufferStats *BufferStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("count", bufferStats.Count, ch)
	group.Collect("max_size_bytes", bufferStats.MaxSizeBytes, ch)
	group.Collect("size_bytes", bufferStats.SizeBytes, ch)
}

type MetricsNetworkStats struct {
	Bytes          float64         `bson:"bytes"`
	Ops            float64         `bson:"ops"`
	GetMores       *BenchmarkStats `bson:"getmores"`
	ReadersCreated float64         `bson:"readersCreated"`
}

func (metricsNetworkStats *MetricsNetworkStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)
	group.Collect("bytes", metricsNetworkStats.Bytes, ch)
	group.Collect("ops", metricsNetworkStats.Ops, ch)
	group.Collect("readers_created", metricsNetworkStats.ReadersCreated, ch)

	metricsNetworkStats.GetMores.Collect(groupName+"_getmores", ch)
}

type ReplStats struct {
	Apply   *ApplyStats          `bson:"apply"`
	Buffer  *BufferStats         `bson:"buffer"`
	Network *MetricsNetworkStats `bson:"network"`
	PreloadStats *PreloadStats   `bson:"preload"`
}

func (replStats *ReplStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	replStats.Apply.Collect(groupName+"_apply", ch)
	replStats.Buffer.Collect(groupName+"_buffer", ch)
	replStats.Network.Collect(groupName+"_network", ch)
	replStats.PreloadStats.Collect(groupName+"_preload", ch)
}

type PreloadStats struct {
	Docs     *BenchmarkStats `bson:"docs"`
	Indexes *BenchmarkStats `bson:"indexes"`
}

func (preloadStats *PreloadStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	preloadStats.Docs.Collect(groupName+"_docs", ch)
	preloadStats.Indexes.Collect(groupName+"_indexes", ch)
}

type StorageStats struct {
	BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
	Requests        float64 `bson:"freelist.search.requests"`
	Scanned         float64 `bson:"freelist.search.scanned"`
}

func (storageStats *StorageStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	group := shared.FindOrCreateGroup(groupName)

	group.Collect("freelist_search_bucket_exhausted", storageStats.BucketExhausted, ch)
	group.Collect("freelist_search_requests", storageStats.Requests, ch)
	group.Collect("freelist_search_scanned", storageStats.Scanned, ch)
}

type MetricsStats struct {
	Document      *DocumentStats      `bson:"document"`
	GetLastError  *GetLastErrorStats  `bson:"getLastError"`
	Operation     *OperationStats     `bson:"operation"`
	QueryExecutor *QueryExecutorStats `bson:"queryExecutor"`
	Record        *RecordStats        `bson:"record"`
	Repl          *ReplStats          `bson:"repl"`
	Storage       *StorageStats       `bson:"storage"`
}

func (metricsStats *MetricsStats) Collect(groupName string, ch chan<- prometheus.Metric) {
	metricsStats.Document.Collect(groupName+"_document", ch)
	metricsStats.GetLastError.Collect(groupName+"_get_last_error", ch)
	metricsStats.Operation.Collect(groupName+"_operation", ch)
	metricsStats.QueryExecutor.Collect(groupName+"_query_executor", ch)
	metricsStats.Record.Collect(groupName+"_record", ch)
	metricsStats.Repl.Collect(groupName+"_repl", ch)
	metricsStats.Storage.Collect(groupName+"_storage", ch)
}
