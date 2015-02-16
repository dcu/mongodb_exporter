package collector

import(
    "testing"
	"github.com/dcu/mongodb_exporter/shared"
)

func Test_MetricsCollectData(t *testing.T) {
	stats := &MetricsStats{
		Document: &DocumentStats{},
		GetLastError: &GetLastErrorStats{
			Wtime: &BenchmarkStats{},
		},
		Operation: &OperationStats{},
		QueryExecutor: &QueryExecutorStats{},
		Record: &RecordStats{},
		Repl: &ReplStats{
			Apply: &ApplyStats{
				Batches: &BenchmarkStats{},
			},
			Buffer: &BufferStats{},
			Network: &MetricsNetworkStats{
				GetMores: &BenchmarkStats{},
			},
			PreloadStats: &PreloadStats{
				Docs: &BenchmarkStats{},
				Indexes: &BenchmarkStats{},
			},
		},
		Storage: &StorageStats{},
	}

	groupName := "metrics"
	stats.Collect(groupName, nil)

	if shared.Groups[groupName+"_document"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_get_last_error"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_operation"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_query_executor"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_record"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_repl_apply"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_storage_freelist_search"] == nil {
		t.Error("Group not created")
	}
}

