package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsCursorTimedOutTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_cursor",
		Name:      "timed_out_total",
		Help:      "timedOut provides the total number of cursors that have timed out since the server process started. If this number is large or growing at a regular rate, this may indicate an application error",
	})
	metricsCursorOpen = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      "metrics_cursor_open",
		Help:      "The open is an embedded document that contains data regarding open cursors",
	}, []string{"state"})
	metricsDocumentTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "metrics_document_total",
		Help:      "The document holds a document of that reflect document access and modification patterns and data use. Compare these values to the data in the opcounters document, which track total number of operations",
	}, []string{"state"})
	metricsGetLastErrorWtimeNumTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "metrics_get_last_error_wtime",
		Name:      "num_total",
		Help:      "num reports the total number of getLastError operations with a specified write concern (i.e. w) that wait for one or more members of a replica set to acknowledge the write operation (i.e. a w value greater than 1.)",
	})
	metricsGetLastErrorWtimeTotalMilliseconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_get_last_error_wtime",
		Name:      "total_milliseconds",
		Help:      "total_millis reports the total amount of time in milliseconds that the mongod has spent performing getLastError operations with write concern (i.e. w) that wait for one or more members of a replica set to acknowledge the write operation (i.e. a w value greater than 1.)",
	})
	metricsGetLastErrorWtimeoutsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_get_last_error",
		Name:      "wtimeouts_total",
		Help:      "wtimeouts reports the number of times that write concern operations have timed out as a result of the wtimeout threshold to getLastError.",
	})
	metricsOperationTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "metrics_operation_total",
		Help:      "operation is a sub-document that holds counters for several types of update and query operations that MongoDB handles using special operation types",
	}, []string{"type"})
	metricsQueryExecutorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "metrics_query_executor_total",
		Help:      "queryExecutor is a document that reports data from the query execution system",
	}, []string{"state"})
	metricsRecordMovesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_record",
		Name:      "moves_total",
		Help:      "moves reports the total number of times documents move within the on-disk representation of the MongoDB data set. Documents move as a result of operations that increase the size of the document beyond their allocated record size",
	})
	metricsReplApplyBatchesNumTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_apply_batches",
		Name:      "num_total",
		Help:      "num reports the total number of batches applied across all databases",
	})
	metricsReplApplyBatchesTotalMilliseconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_apply_batches",
		Name:      "total_milliseconds",
		Help:      "total_millis reports the total amount of time the mongod has spent applying operations from the oplog",
	})
	metricsReplApplyOpsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_apply",
		Name:      "ops_total",
		Help:      "ops reports the total number of oplog operations applied",
	})
	metricsReplBufferCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_buffer",
		Name:      "count",
		Help:      "count reports the current number of operations in the oplog buffer",
	})
	metricsReplBufferMaxSizeBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_buffer",
		Name:      "max_size_bytes",
		Help:      "maxSizeBytes reports the maximum size of the buffer. This value is a constant setting in the mongod, and is not configurable",
	})
	metricsReplBufferSizeBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_buffer",
		Name:      "size_bytes",
		Help:      "sizeBytes reports the current size of the contents of the oplog buffer",
	})
	metricsReplNetworkGetmoresNumTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_network_getmores",
		Name:      "num_total",
		Help:      "num reports the total number of getmore operations, which are operations that request an additional set of operations from the replication sync source.",
	})
	metricsReplNetworkGetmoresTotalMilliseconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_network_getmores",
		Name:      "total_milliseconds",
		Help:      "total_millis reports the total amount of time required to collect data from getmore operations",
	})
	metricsReplNetworkBytesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_network",
		Name:      "bytes_total",
		Help:      "bytes reports the total amount of data read from the replication sync source",
	})
	metricsReplNetworkOpsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_network",
		Name:      "ops_total",
		Help:      "ops reports the total number of operations read from the replication source.",
	})
	metricsReplNetworkReadersCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_network",
		Name:      "readers_created_total",
		Help:      "readersCreated reports the total number of oplog query processes created. MongoDB will create a new oplog query any time an error occurs in the connection, including a timeout, or a network operation. Furthermore, readersCreated will increment every time MongoDB selects a new source fore replication.",
	})
	metricsReplOplogInsertNumTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_oplog_insert",
		Name:      "num_total",
		Help:      "num reports the total number of items inserted into the oplog.",
	})
	metricsReplOplogInsertTotalMilliseconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_oplog_insert",
		Name:      "total_milliseconds",
		Help:      "total_millis reports the total amount of time spent for the mongod to insert data into the oplog.",
	})
	metricsReplOplogInsertBytesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_oplog",
		Name:      "insert_bytes_total",
		Help:      "insertBytes the total size of documents inserted into the oplog.",
	})
	metricsReplPreloadDocsNumTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_preload_docs",
		Name:      "num_total",
		Help:      "num reports the total number of documents loaded during the pre-fetch stage of replication",
	})
	metricsReplPreloadDocsTotalMilliseconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_preload_docs",
		Name:      "total_milliseconds",
		Help:      "total_millis reports the total amount of time spent loading documents as part of the pre-fetch stage of replication",
	})
	metricsReplPreloadIndexesNumTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_preload_indexes",
		Name:      "num_total",
		Help:      "num reports the total number of index entries loaded by members before updating documents as part of the pre-fetch stage of replication",
	})
	metricsReplPreloadIndexesTotalMilliseconds = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_repl_preload_indexes",
		Name:      "total_milliseconds",
		Help:      "total_millis reports the total amount of time spent loading index entries as part of the pre-fetch stage of replication",
	})
	metricsStorageFreelistSearchTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: Namespace,
		Name:      "metrics_storage_freelist_search_total",
		Help:      "metrics about searching records in the database.",
	}, []string{"type"})
	metricsTTLDeletedDocumentsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_ttl",
		Name:      "deleted_documents_total",
		Help:      "deletedDocuments reports the total number of documents deleted from collections with a ttl index.",
	})
	metricsTTLPassesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: "metrics_ttl",
		Name:      "passes_total",
		Help:      "passes reports the number of times the background process removes documents from collections with a ttl index",
	})
)

// DocumentStats are the stats associated to a document.
type DocumentStats struct {
	Deleted  float64 `bson:"deleted"`
	Inserted float64 `bson:"inserted"`
	Returned float64 `bson:"returned"`
	Updated  float64 `bson:"updated"`
}

// Export exposes the document stats to be consumed by the prometheus server.
func (documentStats *DocumentStats) Export(ch chan<- prometheus.Metric) {
	metricsDocumentTotal.WithLabelValues("deleted").Add(documentStats.Deleted)
	metricsDocumentTotal.WithLabelValues("inserted").Add(documentStats.Inserted)
	metricsDocumentTotal.WithLabelValues("returned").Add(documentStats.Returned)
	metricsDocumentTotal.WithLabelValues("updated").Add(documentStats.Updated)
}

// BenchmarkStats is bechmark info about an operation.
type BenchmarkStats struct {
	Num         float64 `bson:"num"`
	TotalMillis float64 `bson:"totalMillis"`
}

// GetLastErrorStats are the last error stats.
type GetLastErrorStats struct {
	Wtimeouts float64         `bson:"wtimeouts"`
	Wtime     *BenchmarkStats `bson:"wtime"`
}

// Export exposes the get last error stats.
func (getLastErrorStats *GetLastErrorStats) Export(ch chan<- prometheus.Metric) {
	metricsGetLastErrorWtimeNumTotal.Set(getLastErrorStats.Wtime.Num)
	metricsGetLastErrorWtimeTotalMilliseconds.Add(getLastErrorStats.Wtime.TotalMillis)

	metricsGetLastErrorWtimeoutsTotal.Add(getLastErrorStats.Wtimeouts)
}

// OperationStats are the stats for some kind of operations.
type OperationStats struct {
	Fastmod      float64 `bson:"fastmod"`
	Idhack       float64 `bson:"idhack"`
	ScanAndOrder float64 `bson:"scanAndOrder"`
}

// Export exports the operation stats.
func (operationStats *OperationStats) Export(ch chan<- prometheus.Metric) {
	metricsOperationTotal.WithLabelValues("fastmod").Add(operationStats.Fastmod)
	metricsOperationTotal.WithLabelValues("idhack").Add(operationStats.Idhack)
	metricsOperationTotal.WithLabelValues("scan_and_order").Add(operationStats.ScanAndOrder)
}

// QueryExecutorStats are the stats associated with a query execution.
type QueryExecutorStats struct {
	Scanned        float64 `bson:"scanned"`
	ScannedObjects float64 `bson:"scannedObjects"`
}

// Export exports the query executor stats.
func (queryExecutorStats *QueryExecutorStats) Export(ch chan<- prometheus.Metric) {
	metricsQueryExecutorTotal.WithLabelValues("scanned").Add(queryExecutorStats.Scanned)
	metricsQueryExecutorTotal.WithLabelValues("scanned_objects").Add(queryExecutorStats.ScannedObjects)
}

// RecordStats are stats associated with a record.
type RecordStats struct {
	Moves float64 `bson:"moves"`
}

// Export exposes the record stats.
func (recordStats *RecordStats) Export(ch chan<- prometheus.Metric) {
	metricsRecordMovesTotal.Add(recordStats.Moves)
}

// ApplyStats are the stats associated with the apply operation.
type ApplyStats struct {
	Batches *BenchmarkStats `bson:"batches"`
	Ops     float64         `bson:"ops"`
}

// Export exports the apply stats
func (applyStats *ApplyStats) Export(ch chan<- prometheus.Metric) {
	metricsReplApplyOpsTotal.Add(applyStats.Ops)

	metricsReplApplyBatchesNumTotal.Add(applyStats.Batches.Num)
	metricsReplApplyBatchesTotalMilliseconds.Add(applyStats.Batches.TotalMillis)
}

// BufferStats are the stats associated with the buffer
type BufferStats struct {
	Count        float64 `bson:"count"`
	MaxSizeBytes float64 `bson:"maxSizeBytes"`
	SizeBytes    float64 `bson:"sizeBytes"`
}

// Export exports the buffer stats.
func (bufferStats *BufferStats) Export(ch chan<- prometheus.Metric) {
	metricsReplBufferCount.Set(bufferStats.Count)
	metricsReplBufferMaxSizeBytes.Add(bufferStats.MaxSizeBytes)
	metricsReplBufferSizeBytes.Set(bufferStats.SizeBytes)
}

// MetricsNetworkStats are the network stats.
type MetricsNetworkStats struct {
	Bytes          float64         `bson:"bytes"`
	Ops            float64         `bson:"ops"`
	GetMores       *BenchmarkStats `bson:"getmores"`
	ReadersCreated float64         `bson:"readersCreated"`
}

// Export exposes the network stats.
func (metricsNetworkStats *MetricsNetworkStats) Export(ch chan<- prometheus.Metric) {
	metricsReplNetworkBytesTotal.Add(metricsNetworkStats.Bytes)
	metricsReplNetworkOpsTotal.Add(metricsNetworkStats.Ops)
	metricsReplNetworkReadersCreatedTotal.Add(metricsNetworkStats.ReadersCreated)

	metricsReplNetworkGetmoresNumTotal.Add(metricsNetworkStats.GetMores.Num)
	metricsReplNetworkGetmoresTotalMilliseconds.Add(metricsNetworkStats.GetMores.TotalMillis)
}

// ReplStats are the stats associated with the replication process.
type ReplStats struct {
	Apply        *ApplyStats          `bson:"apply"`
	Buffer       *BufferStats         `bson:"buffer"`
	Network      *MetricsNetworkStats `bson:"network"`
	PreloadStats *PreloadStats        `bson:"preload"`
}

// Export exposes the replication stats.
func (replStats *ReplStats) Export(ch chan<- prometheus.Metric) {
	replStats.Apply.Export(ch)
	replStats.Buffer.Export(ch)
	replStats.Network.Export(ch)
	replStats.PreloadStats.Export(ch)
}

// PreloadStats are the stats associated with preload operation.
type PreloadStats struct {
	Docs    *BenchmarkStats `bson:"docs"`
	Indexes *BenchmarkStats `bson:"indexes"`
}

// Export exposes the preload stats.
func (preloadStats *PreloadStats) Export(ch chan<- prometheus.Metric) {
	metricsReplPreloadDocsNumTotal.Add(preloadStats.Docs.Num)
	metricsReplPreloadDocsTotalMilliseconds.Add(preloadStats.Docs.TotalMillis)

	metricsReplPreloadIndexesNumTotal.Add(preloadStats.Indexes.Num)
	metricsReplPreloadIndexesTotalMilliseconds.Add(preloadStats.Indexes.TotalMillis)
}

// StorageStats are the stats associated with the storage.
type StorageStats struct {
	BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
	Requests        float64 `bson:"freelist.search.requests"`
	Scanned         float64 `bson:"freelist.search.scanned"`
}

// Export exports the storage stats.
func (storageStats *StorageStats) Export(ch chan<- prometheus.Metric) {
	metricsStorageFreelistSearchTotal.WithLabelValues("bucket_exhausted").Add(storageStats.BucketExhausted)
	metricsStorageFreelistSearchTotal.WithLabelValues("requests").Add(storageStats.Requests)
	metricsStorageFreelistSearchTotal.WithLabelValues("scanned").Add(storageStats.Scanned)
}

// CursorStatsOpen are the stats for open cursors
type CursorStatsOpen struct {
	NoTimeout float64 `bson:"noTimeout"`
	Pinned    float64 `bson:"pinned"`
	Total     float64 `bson:"total"`
}

// CursorStats are the stats for cursors
type CursorStats struct {
	TimedOut float64          `bson:"timedOut"`
	Open     *CursorStatsOpen `bson:"open"`
}

// Export exports the cursor stats.
func (cursorStats *CursorStats) Export(ch chan<- prometheus.Metric) {
	metricsCursorTimedOutTotal.Add(cursorStats.TimedOut)
	metricsCursorOpen.WithLabelValues("timed_out").Set(cursorStats.Open.NoTimeout)
	metricsCursorOpen.WithLabelValues("pinned").Set(cursorStats.Open.Pinned)
	metricsCursorOpen.WithLabelValues("total").Set(cursorStats.Open.Total)
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
	Cursor        *CursorStats        `bson:"cursor"`
}

// Export exports the metrics stats.
func (metricsStats *MetricsStats) Export(ch chan<- prometheus.Metric) {
	if metricsStats.Document != nil {
		metricsStats.Document.Export(ch)
	}
	if metricsStats.GetLastError != nil {
		metricsStats.GetLastError.Export(ch)
	}
	if metricsStats.Operation != nil {
		metricsStats.Operation.Export(ch)
	}
	if metricsStats.QueryExecutor != nil {
		metricsStats.QueryExecutor.Export(ch)
	}
	if metricsStats.Record != nil {
		metricsStats.Record.Export(ch)
	}
	if metricsStats.Repl != nil {
		metricsStats.Repl.Export(ch)
	}
	if metricsStats.Storage != nil {
		metricsStats.Storage.Export(ch)
	}
	if metricsStats.Cursor != nil {
		metricsStats.Cursor.Export(ch)
	}

	metricsCursorTimedOutTotal.Collect(ch)
	metricsCursorOpen.Collect(ch)
	metricsDocumentTotal.Collect(ch)
	metricsGetLastErrorWtimeNumTotal.Collect(ch)
	metricsGetLastErrorWtimeTotalMilliseconds.Collect(ch)
	metricsGetLastErrorWtimeoutsTotal.Collect(ch)
	metricsOperationTotal.Collect(ch)
	metricsQueryExecutorTotal.Collect(ch)
	metricsRecordMovesTotal.Collect(ch)
	metricsReplApplyBatchesNumTotal.Collect(ch)
	metricsReplApplyBatchesTotalMilliseconds.Collect(ch)
	metricsReplApplyOpsTotal.Collect(ch)
	metricsReplBufferCount.Collect(ch)
	metricsReplBufferMaxSizeBytes.Collect(ch)
	metricsReplBufferSizeBytes.Collect(ch)
	metricsReplNetworkGetmoresNumTotal.Collect(ch)
	metricsReplNetworkGetmoresTotalMilliseconds.Collect(ch)
	metricsReplNetworkBytesTotal.Collect(ch)
	metricsReplNetworkOpsTotal.Collect(ch)
	metricsReplNetworkReadersCreatedTotal.Collect(ch)
	metricsReplOplogInsertNumTotal.Collect(ch)
	metricsReplOplogInsertTotalMilliseconds.Collect(ch)
	metricsReplOplogInsertBytesTotal.Collect(ch)
	metricsReplPreloadDocsNumTotal.Collect(ch)
	metricsReplPreloadDocsTotalMilliseconds.Collect(ch)
	metricsReplPreloadIndexesNumTotal.Collect(ch)
	metricsReplPreloadIndexesTotalMilliseconds.Collect(ch)
	metricsStorageFreelistSearchTotal.Collect(ch)
	metricsTTLDeletedDocumentsTotal.Collect(ch)
	metricsTTLPassesTotal.Collect(ch)
}

// Describe describes the metrics for prometheus
func (metricsStats *MetricsStats) Describe(ch chan<- *prometheus.Desc) {
	metricsCursorTimedOutTotal.Describe(ch)
	metricsCursorOpen.Describe(ch)
	metricsDocumentTotal.Describe(ch)
	metricsGetLastErrorWtimeNumTotal.Describe(ch)
	metricsGetLastErrorWtimeTotalMilliseconds.Describe(ch)
	metricsGetLastErrorWtimeoutsTotal.Describe(ch)
	metricsOperationTotal.Describe(ch)
	metricsQueryExecutorTotal.Describe(ch)
	metricsRecordMovesTotal.Describe(ch)
	metricsReplApplyBatchesNumTotal.Describe(ch)
	metricsReplApplyBatchesTotalMilliseconds.Describe(ch)
	metricsReplApplyOpsTotal.Describe(ch)
	metricsReplBufferCount.Describe(ch)
	metricsReplBufferMaxSizeBytes.Describe(ch)
	metricsReplBufferSizeBytes.Describe(ch)
	metricsReplNetworkGetmoresNumTotal.Describe(ch)
	metricsReplNetworkGetmoresTotalMilliseconds.Describe(ch)
	metricsReplNetworkBytesTotal.Describe(ch)
	metricsReplNetworkOpsTotal.Describe(ch)
	metricsReplNetworkReadersCreatedTotal.Describe(ch)
	metricsReplOplogInsertNumTotal.Describe(ch)
	metricsReplOplogInsertTotalMilliseconds.Describe(ch)
	metricsReplOplogInsertBytesTotal.Describe(ch)
	metricsReplPreloadDocsNumTotal.Describe(ch)
	metricsReplPreloadDocsTotalMilliseconds.Describe(ch)
	metricsReplPreloadIndexesNumTotal.Describe(ch)
	metricsReplPreloadIndexesTotalMilliseconds.Describe(ch)
	metricsStorageFreelistSearchTotal.Describe(ch)
	metricsTTLDeletedDocumentsTotal.Describe(ch)
	metricsTTLPassesTotal.Describe(ch)
}
