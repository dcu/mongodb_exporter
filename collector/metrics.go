package collector

import(
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/shared"
)

// Metrics
type DocumentStats struct {
    Deleted float64 `bson:"deleted" type:"gauge"`
    Inserted float64 `bson:"inserted" type:"gauge"`
    Returned float64 `bson:"returned" type:"gauge"`
    Updated float64 `bson:"updated" type:"gauge"`
}

func (documentStats *DocumentStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)

    group.Collect("deleted", documentStats.Deleted, ch)
    group.Collect("inserted", documentStats.Inserted, ch)
    group.Collect("returned", documentStats.Returned, ch)
    group.Collect("updated", documentStats.Updated, ch)
}

type BenchmarkStats struct {
    Num float64 `bson:"num" type:"counter"`
    TotalMillis float64 `bson:"totalMillis" type:"counter"`
}

func (benchmarkStats *BenchmarkStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)

    group.Collect("num", benchmarkStats.Num, ch)
    group.Collect("total_millis", benchmarkStats.TotalMillis, ch)
}

type GetLastErrorStats struct {
    Wtimeouts float64 `bson:"wtimeouts" type:"counter"`
    Wtime *BenchmarkStats `bson:"wtime" type:"counter"`
}

func (getLastErrorStats *GetLastErrorStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)

    group.Collect("wtimeouts", getLastErrorStats.Wtimeouts, ch)
    getLastErrorStats.Wtime.Collect(groupName+"_wtime", ch)
}

type OperationStats struct {
    Fastmod float64 `bson:"fastmod" type:"counter"`
    Idhack float64 `bson:"idhack" type:"counter"`
    ScanAndOrder float64 `bson:"scanAndOrder" type:"counter"`
}

func (operationStats *OperationStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("fastmod", operationStats.Fastmod, ch)
    group.Collect("idhack", operationStats.Idhack, ch)
    group.Collect("scan_and_order", operationStats.ScanAndOrder, ch)
}

type QueryExecutorStats struct {
    Scanned float64 `bson:"scanned" type:"counter"`
    ScannedObjects float64 `bson:"scannedObjects" type:"counter"`
}

func (queryExecutorStats *QueryExecutorStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("scanned", queryExecutorStats.Scanned, ch)
    group.Collect("scanned_objects", queryExecutorStats.ScannedObjects, ch)
}

type RecordStats struct {
    Moves float64 `bson:"moves" type:"counter"`
}

func (recordStats *RecordStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("moves", recordStats.Moves, ch)
}

type ApplyStats struct {
    Batches *BenchmarkStats `bson:"batches"`
    Ops float64 `bson:"ops" type:"counter"`
}

func (applyStats *ApplyStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("ops", applyStats.Ops, ch)

    applyStats.Batches.Collect(groupName+"_batches", ch)
}

type BufferStats struct {
    Count float64 `bson:"count" type:"counter"`
    MaxSizeBytes float64 `bson:"maxSizeBytes" type:"counter"`
    SizeBytes float64 `bson:"sizeBytes" type:"counter"`
}

func (bufferStats *BufferStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("count", bufferStats.Count, ch)
    group.Collect("max_size_bytes", bufferStats.MaxSizeBytes, ch)
    group.Collect("size_bytes", bufferStats.SizeBytes, ch)
}

type MetricsNetworkStats struct {
    Bytes float64 `bson:"bytes" type:"counter"`
    Ops float64 `bson:"ops" type:"counter"`
    GetMores *BenchmarkStats `bson:"getmores"`
    ReadersCreated float64 `bson:"readersCreated"`
}

func (metricsNetworkStats *MetricsNetworkStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)
    group.Collect("bytes", metricsNetworkStats.Bytes, ch)
    group.Collect("ops", metricsNetworkStats.Ops, ch)
    group.Collect("readers_created", metricsNetworkStats.ReadersCreated, ch)

    metricsNetworkStats.GetMores.Collect(groupName+"_getmores", ch)
}

type ReplStats struct {
    Apply *ApplyStats `bson:"apply" group:"document"`
    Buffer *BufferStats `bson:"buffer" group:"document"`
    Network *MetricsNetworkStats `bson:"network" group:"document"`
}

func (replStats *ReplStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    replStats.Apply.Collect(groupName+"_apply", ch)
    replStats.Buffer.Collect(groupName+"_buffer", ch)
    replStats.Network.Collect(groupName+"_network", ch)
}

type PreloadStats struct {
    Doc *BenchmarkStats `bson:"doc"`
    Indexes *BenchmarkStats `bson:"indexes"`
}

func (preloadStats *PreloadStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    preloadStats.Doc.Collect(groupName+"_doc", ch)
    preloadStats.Indexes.Collect(groupName+"_indexes", ch)
}

type StorageStats struct {
    BucketExhausted float64 `bson:"freelist.search.bucketExhausted"`
    Requests float64 `bson:"freelist.search.requests"`
    Scanned float64 `bson:"freelist.search.scanned"`
}

func (storageStats *StorageStats) Collect(groupName string, ch chan<- prometheus.Metric) {
    group := shared.FindOrCreateGroup(groupName)

    group.Collect("freelist_search_bucket_exhausted", storageStats.BucketExhausted, ch)
    group.Collect("freelist_search_requests", storageStats.Requests, ch)
    group.Collect("freelist_search_scanned", storageStats.Scanned, ch)
}

type MetricsStats struct {
    Document *DocumentStats `bson:"document" group:"document"`
    GetLastError *GetLastErrorStats `bson:"getLastError" group:"document"`
    Operation *OperationStats `bson:"operation" group:"document"`
    QueryExecutor *QueryExecutorStats `bson:"queryExecutor" group:"document"`
    Record *RecordStats `bson:"record" group:"document"`
    Repl *ReplStats `bson:"repl" group:"document"`
    Storage *StorageStats `bson:"storage" group:"document"`
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
