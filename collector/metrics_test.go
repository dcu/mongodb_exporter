package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"testing"
)

func Test_MetricsCollectData(t *testing.T) {
	stats := &MetricsStats{
		Document: &DocumentStats{},
		GetLastError: &GetLastErrorStats{
			Wtime: &BenchmarkStats{},
		},
		Operation:     &OperationStats{},
		QueryExecutor: &QueryExecutorStats{},
		Record:        &RecordStats{},
		Repl: &ReplStats{
			Apply: &ApplyStats{
				Batches: &BenchmarkStats{},
			},
			Buffer: &BufferStats{},
			Network: &MetricsNetworkStats{
				GetMores: &BenchmarkStats{},
			},
			PreloadStats: &PreloadStats{
				Docs:    &BenchmarkStats{},
				Indexes: &BenchmarkStats{},
			},
		},
		Storage: &StorageStats{},
	}

	groupName := "metrics"
	stats.Export(groupName)

	if shared.Groups[groupName+"_document_total"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_get_last_error"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_operation_total"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_query_executor_total"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_record"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_repl_apply"] == nil {
		t.Error("Group not created")
	}
	if shared.Groups[groupName+"_storage_freelist_search_total"] == nil {
		t.Error("Group not created")
	}
}
