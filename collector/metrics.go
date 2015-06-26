package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

// DocumentStats are the stats associated to a document.
type DocumentStats struct {
	Deleted  float64 `bson:"deleted"`
	Inserted float64 `bson:"inserted"`
	Returned float64 `bson:"returned"`
	Updated  float64 `bson:"updated"`
}

// Export exposes the document stats to be consumed by the prometheus server.
func (documentStats *DocumentStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("deleted", documentStats.Deleted)
	group.Export("inserted", documentStats.Inserted)
	group.Export("returned", documentStats.Returned)
	group.Export("updated", documentStats.Updated)
}

// BenchmarkStats is bechmark info about an operation.
type BenchmarkStats struct {
	Num         float64 `bson:"num"`
	TotalMillis float64 `bson:"totalMillis"`
}

// Export exports the benchmark stats.
func (benchmarkStats *BenchmarkStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("num_total", benchmarkStats.Num)
	group.Export("total_milliseconds", benchmarkStats.TotalMillis)
}

// GetLastErrorStats are the last error stats.
type GetLastErrorStats struct {
	Wtimeouts float64         `bson:"wtimeouts"`
	Wtime     *BenchmarkStats `bson:"wtime"`
}

// Export exposes the get last error stats.
func (getLastErrorStats *GetLastErrorStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("wtimeouts_total", getLastErrorStats.Wtimeouts)
	getLastErrorStats.Wtime.Export(groupName + "_wtime")
}

// OperationStats are the stats for some kind of operations.
type OperationStats struct {
	Fastmod      float64 `bson:"fastmod"`
	Idhack       float64 `bson:"idhack"`
	ScanAndOrder float64 `bson:"scanAndOrder"`
}

// Export exports the operation stats.
func (operationStats *OperationStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("fastmod", operationStats.Fastmod)
	group.Export("idhack", operationStats.Idhack)
	group.Export("scan_and_order", operationStats.ScanAndOrder)
}

// QueryExecutorStats are the stats associated with a query execution.
type QueryExecutorStats struct {
	Scanned        float64 `bson:"scanned"`
	ScannedObjects float64 `bson:"scannedObjects"`
}

// Export exports the query executor stats.
func (queryExecutorStats *QueryExecutorStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("scanned", queryExecutorStats.Scanned)
	group.Export("scanned_objects", queryExecutorStats.ScannedObjects)
}

// RecordStats are stats associated with a record.
type RecordStats struct {
	Moves float64 `bson:"moves"`
}

// Export exposes the record stats.
func (recordStats *RecordStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("moves_total", recordStats.Moves)
}

// ApplyStats are the stats associated with the apply operation.
type ApplyStats struct {
	Batches *BenchmarkStats `bson:"batches"`
	Ops     float64         `bson:"ops"`
}

// Export exports the apply stats
func (applyStats *ApplyStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("ops_total", applyStats.Ops)

	applyStats.Batches.Export(groupName + "_batches")
}

// BufferStats are the stats associated with the buffer
type BufferStats struct {
	Count        float64 `bson:"count"`
	MaxSizeBytes float64 `bson:"maxSizeBytes"`
	SizeBytes    float64 `bson:"sizeBytes"`
}

// Export exports the buffer stats.
func (bufferStats *BufferStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("count", bufferStats.Count)
	group.Export("max_size_bytes", bufferStats.MaxSizeBytes)
	group.Export("size_bytes", bufferStats.SizeBytes)
}

// MetricsNetworkStats are the network stats.
type MetricsNetworkStats struct {
	Bytes          float64         `bson:"bytes"`
	Ops            float64         `bson:"ops"`
	GetMores       *BenchmarkStats `bson:"getmores"`
	ReadersCreated float64         `bson:"readersCreated"`
}

// Export exposes the network stats.
func (metricsNetworkStats *MetricsNetworkStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("bytes_total", metricsNetworkStats.Bytes)
	group.Export("ops_total", metricsNetworkStats.Ops)
	group.Export("readers_created_total", metricsNetworkStats.ReadersCreated)

	metricsNetworkStats.GetMores.Export(groupName + "_getmores")
}

// ReplStats are the stats associated with the replication process.
type ReplStats struct {
	Apply        *ApplyStats          `bson:"apply"`
	Buffer       *BufferStats         `bson:"buffer"`
	Network      *MetricsNetworkStats `bson:"network"`
	PreloadStats *PreloadStats        `bson:"preload"`
}

// Export exposes the replication stats.
func (replStats *ReplStats) Export(groupName string) {
	replStats.Apply.Export(groupName + "_apply")
	replStats.Buffer.Export(groupName + "_buffer")
	replStats.Network.Export(groupName + "_network")
	replStats.PreloadStats.Export(groupName + "_preload")
}

// PreloadStats are the stats associated with preload operation.
type PreloadStats struct {
	Docs    *BenchmarkStats `bson:"docs"`
	Indexes *BenchmarkStats `bson:"indexes"`
}

// Export exposes the preload stats.
func (preloadStats *PreloadStats) Export(groupName string) {
	preloadStats.Docs.Export(groupName + "_docs")
	preloadStats.Indexes.Export(groupName + "_indexes")
}

// StorageStats are the stats associated with the storage.
type StorageStats struct {
	BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
	Requests        float64 `bson:"freelist.search.requests"`
	Scanned         float64 `bson:"freelist.search.scanned"`
}

// Export exports the storage stats.
func (storageStats *StorageStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("bucket_exhausted", storageStats.BucketExhausted)
	group.Export("requests", storageStats.Requests)
	group.Export("scanned", storageStats.Scanned)
}

// MetricsStats are all stats associated with metrics of the system
type MetricsStats struct {
	Document      *DocumentStats      `bson:"document"`
	GetLastError  *GetLastErrorStats  `bson:"getLastError"`
	Operation     *OperationStats     `bson:"operation"`
	QueryExecutor *QueryExecutorStats `bson:"queryExecutor"`
	Record        *RecordStats        `bson:"record"`
	Repl          *ReplStats          `bson:"repl"`
	Storage       *StorageStats       `bson:"storage"`
}

// Export exports the metrics stats.
func (metricsStats *MetricsStats) Export(groupName string) {
	if metricsStats.Document != nil {
		metricsStats.Document.Export(groupName + "_document_total")
	}
	if metricsStats.GetLastError != nil {
		metricsStats.GetLastError.Export(groupName + "_get_last_error")
	}
	if metricsStats.Operation != nil {
		metricsStats.Operation.Export(groupName + "_operation_total")
	}
	if metricsStats.QueryExecutor != nil {
		metricsStats.QueryExecutor.Export(groupName + "_query_executor_total")
	}
	if metricsStats.Record != nil {
		metricsStats.Record.Export(groupName + "_record")
	}
	if metricsStats.Repl != nil {
		metricsStats.Repl.Export(groupName + "_repl")
	}
	if metricsStats.Storage != nil {
		metricsStats.Storage.Export(groupName + "_storage_freelist_search_total")
	}
}
