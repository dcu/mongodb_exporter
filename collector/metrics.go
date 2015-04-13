package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
)

// Metrics
type DocumentStats struct {
	Deleted  float64 `bson:"deleted"`
	Inserted float64 `bson:"inserted"`
	Returned float64 `bson:"returned"`
	Updated  float64 `bson:"updated"`
}

func (documentStats *DocumentStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("deleted", documentStats.Deleted)
	group.Export("inserted", documentStats.Inserted)
	group.Export("returned", documentStats.Returned)
	group.Export("updated", documentStats.Updated)
}

type BenchmarkStats struct {
	Num         float64 `bson:"num"`
	TotalMillis float64 `bson:"totalMillis"`
}

func (benchmarkStats *BenchmarkStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("num_total", benchmarkStats.Num)
	group.Export("total_milliseconds", benchmarkStats.TotalMillis)
}

type GetLastErrorStats struct {
	Wtimeouts float64         `bson:"wtimeouts"`
	Wtime     *BenchmarkStats `bson:"wtime"`
}

func (getLastErrorStats *GetLastErrorStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("wtimeouts_total", getLastErrorStats.Wtimeouts)
	getLastErrorStats.Wtime.Export(groupName + "_wtime")
}

type OperationStats struct {
	Fastmod      float64 `bson:"fastmod"`
	Idhack       float64 `bson:"idhack"`
	ScanAndOrder float64 `bson:"scanAndOrder"`
}

func (operationStats *OperationStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("fastmod", operationStats.Fastmod)
	group.Export("idhack", operationStats.Idhack)
	group.Export("scan_and_order", operationStats.ScanAndOrder)
}

type QueryExecutorStats struct {
	Scanned        float64 `bson:"scanned"`
	ScannedObjects float64 `bson:"scannedObjects"`
}

func (queryExecutorStats *QueryExecutorStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("scanned", queryExecutorStats.Scanned)
	group.Export("scanned_objects", queryExecutorStats.ScannedObjects)
}

type RecordStats struct {
	Moves float64 `bson:"moves"`
}

func (recordStats *RecordStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("moves_total", recordStats.Moves)
}

type ApplyStats struct {
	Batches *BenchmarkStats `bson:"batches"`
	Ops     float64         `bson:"ops"`
}

func (applyStats *ApplyStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("ops_total", applyStats.Ops)

	applyStats.Batches.Export(groupName + "_batches")
}

type BufferStats struct {
	Count        float64 `bson:"count"`
	MaxSizeBytes float64 `bson:"maxSizeBytes"`
	SizeBytes    float64 `bson:"sizeBytes"`
}

func (bufferStats *BufferStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("count", bufferStats.Count)
	group.Export("max_size_bytes", bufferStats.MaxSizeBytes)
	group.Export("size_bytes", bufferStats.SizeBytes)
}

type MetricsNetworkStats struct {
	Bytes          float64         `bson:"bytes"`
	Ops            float64         `bson:"ops"`
	GetMores       *BenchmarkStats `bson:"getmores"`
	ReadersCreated float64         `bson:"readersCreated"`
}

func (metricsNetworkStats *MetricsNetworkStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)
	group.Export("bytes_total", metricsNetworkStats.Bytes)
	group.Export("ops_total", metricsNetworkStats.Ops)
	group.Export("readers_created_total", metricsNetworkStats.ReadersCreated)

	metricsNetworkStats.GetMores.Export(groupName + "_getmores")
}

type ReplStats struct {
	Apply        *ApplyStats          `bson:"apply"`
	Buffer       *BufferStats         `bson:"buffer"`
	Network      *MetricsNetworkStats `bson:"network"`
	PreloadStats *PreloadStats        `bson:"preload"`
}

func (replStats *ReplStats) Export(groupName string) {
	replStats.Apply.Export(groupName + "_apply")
	replStats.Buffer.Export(groupName + "_buffer")
	replStats.Network.Export(groupName + "_network")
	replStats.PreloadStats.Export(groupName + "_preload")
}

type PreloadStats struct {
	Docs    *BenchmarkStats `bson:"docs"`
	Indexes *BenchmarkStats `bson:"indexes"`
}

func (preloadStats *PreloadStats) Export(groupName string) {
	preloadStats.Docs.Export(groupName + "_docs")
	preloadStats.Indexes.Export(groupName + "_indexes")
}

type StorageStats struct {
	BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
	Requests        float64 `bson:"freelist.search.requests"`
	Scanned         float64 `bson:"freelist.search.scanned"`
}

func (storageStats *StorageStats) Export(groupName string) {
	group := shared.FindOrCreateGroup(groupName)

	group.Export("bucket_exhausted", storageStats.BucketExhausted)
	group.Export("requests", storageStats.Requests)
	group.Export("scanned", storageStats.Scanned)
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

func (metricsStats *MetricsStats) Export(groupName string) {
	metricsStats.Document.Export(groupName + "_document_total")
	metricsStats.GetLastError.Export(groupName + "_get_last_error")
	metricsStats.Operation.Export(groupName + "_operation_total")
	metricsStats.QueryExecutor.Export(groupName + "_query_executor_total")
	metricsStats.Record.Export(groupName + "_record")
	metricsStats.Repl.Export(groupName + "_repl")
	metricsStats.Storage.Export(groupName + "_storage_freelist_search_total")
}
