package shared

import (
	"strings"
)

// FieldDesc is information about the field. It contains the type of the metric, labels to track and a help message.
type FieldDesc struct {
	Type   string
	Labels []string
	Help   string
}

// GroupFieldsMap is a map with string as key and FieldDesc as value
type GroupFieldsMap map[string]*FieldDesc

// GroupDescMap is a map with string as key and GroupFieldsMap as value
type GroupDescMap map[string]GroupFieldsMap

var (
	// GroupsDesc contains all supported groups.
	GroupsDesc = GroupDescMap{
		"instance": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "Information about the server instance.",
			},
			"uptime_seconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The value of the uptime field corresponds to the number of seconds that the mongos or mongod process has been active.",
			},
			"uptime_estimate_seconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "uptimeEstimate provides the uptime as calculated from MongoDB's internal course-grained time keeping system.",
			},
			"local_time": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The localTime value is the current time, according to the server, in UTC specified in an ISODate format.",
			},
		},
		"asserts_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type"},
				Help:   "The asserts document reports the number of asserts on the database. While assert errors are typically uncommon, if there are non-zero values for the asserts, you should check the log file for the mongod process for more information. In many cases these errors are trivial, but are worth investigating.",
			},
			"regular": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The regular counter tracks the number of regular assertions raised since the server process started. Check the log file for more information about these messages.",
			},
			"warning": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The warning counter tracks the number of warnings raised since the server process started. Check the log file for more information about these warnings.",
			},
			"msg": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The msg counter tracks the number of message assertions raised since the server process started. Check the log file for more information about these messages.",
			},
			"user": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The user counter reports the number of \"user asserts\" that have occurred since the last time the server process started. These are errors that user may generate, such as out of disk space or duplicate key. You can prevent these assertions by fixing a problem with your application or deployment. Check the MongoDB log for more information.",
			},
			"rollovers": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The rollovers counter displays the number of times that the rollover counters have rolled over since the last time the server process started. The counters will rollover to zero after 230 assertions. Use this value to provide context to the other values in the asserts data structure.",
			},
		},
		"background_flushing": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "mongod periodically flushes writes to disk. In the default configuration, this happens every 60 seconds. The backgroundFlushing data structure contains data regarding these operations. Consider these values if you have concerns about write performance and journaling",
			},
			"flushes_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "flushes is a counter that collects the number of times the database has flushed all writes to disk. This value will grow as database runs for longer periods of time",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The total_ms value provides the total number of milliseconds (ms) that the mongod processes have spent writing (i.e. flushing) data to disk. Because this is an absolute value, consider the value offlushes and average_ms to provide better context for this datum",
			},
			"average_milliseconds": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The average_ms value describes the relationship between the number of flushes and the total amount of time that the database has spent writing data to disk. The larger flushes is, the more likely this value is likely to represent a \"normal,\" time; however, abnormal data can skew this value",
			},
			"last_milliseconds": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of the last_ms field is the amount of time, in milliseconds, that the last flush operation took to complete. Use this value to verify that the current performance of the server and is in line with the historical data provided by average_ms and total_ms",
			},
			"last_finished_time": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The last_finished field provides a timestamp of the last completed flush operation in the ISODateformat. If this value is more than a few minutes old relative to your server’s current time and accounting for differences in time zone, restarting the database may result in some data loss",
			},
		},
		"connections": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"state"},
				Help:   "The connections sub document data regarding the current status of incoming connections and availability of the database server. Use these values to assess the current load and capacity requirements of the server",
			},
			"current": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The value of current corresponds to the number of connections to the database server from clients. This number includes the current shell session. Consider the value of available to add more context to this datum",
			},
			"available": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "available provides a count of the number of unused available incoming connections the database can provide. Consider this value in combination with the value of current to understand the connection load on the database, and the UNIX ulimit Settings document for more information about system thresholds on available connections",
			},
		},
		"connections_metrics": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "Total connections",
			},
			"created_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "totalCreated provides a count of all incoming connections created to the server. This number includes connections that have since closed",
			},
		},
		"durability_commits": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"state"},
				Help:   "Durability commits",
			},
			"written": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The commits provides the number of transactions written to the journal during the last journal group commit interval.",
			},
			"in_write_lock": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The commitsInWriteLock provides a count of the commits that occurred while a write lock was held. Commits in a write lock indicate a MongoDB node under a heavy write load and call for further diagnosis",
			},
		},
		"durability": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The dur (for “durability”) document contains data regarding the mongod‘s journaling-related operations and performance. mongod must be running with journaling for these data to appear in the output of \"serverStatus\". MongoDB reports the data in dur based on 3 second intervals of data, collected between 3 and 6 seconds in the past",
			},
			"journaled_megabytes": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The journaledMB provides the amount of data in megabytes (MB) written to journal during the last journal group commit interval",
			},
			"write_to_data_files_megabytes": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The writeToDataFilesMB provides the amount of data in megabytes (MB) written from journal to the data files during the last journal group commit interval",
			},
			"compression": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The compression represents the compression ratio of the data written to the journal: ( journaled_size_of_data / uncompressed_size_of_data )",
			},
			"early_commits": &FieldDesc{
				Type:   "summary",
				Labels: []string{},
				Help:   "The earlyCommits value reflects the number of times MongoDB requested a commit before the scheduled journal group commit interval. Use this value to ensure that your journal group commit interval is not too long for your deployment",
			},
		},
		"durability_time_milliseconds": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "summary_vec",
				Labels: []string{"stage"},
				Help:   "Summary of times spent during the journaling process.",
			},
			"dt": &FieldDesc{
				Type:   "summary",
				Labels: []string{},
				Help:   "The dt value provides, in milliseconds, the amount of time over which MongoDB collected the timeMSdata. Use this field to provide context to the other timeMS field values",
			},
			"prep_log_buffer": &FieldDesc{
				Type:   "summary",
				Labels: []string{},
				Help:   "The prepLogBuffer value provides, in milliseconds, the amount of time spent preparing to write to the journal. Smaller values indicate better journal performance",
			},
			"write_to_journal": &FieldDesc{
				Type:   "summary",
				Labels: []string{},
				Help:   "The writeToJournal value provides, in milliseconds, the amount of time spent actually writing to the journal. File system speeds and device interfaces can affect performance",
			},
			"write_to_data_files": &FieldDesc{
				Type:   "summary",
				Labels: []string{},
				Help:   "The writeToDataFiles value provides, in milliseconds, the amount of time spent writing to data files after journaling. File system speeds and device interfaces can affect performance",
			},
			"remap_private_view": &FieldDesc{
				Type:   "summary",
				Labels: []string{},
				Help:   "The remapPrivateView value provides, in milliseconds, the amount of time spent remapping copy-on-write memory mapped views. Smaller values indicate better journal performance",
			},
		},
		"extra_info": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The extra_info data structure holds data collected by the mongod instance about the underlying system. Your system may only report a subset of these fields",
			},
			"page_faults_total": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The page_faults Reports the total number of page faults that require disk operations. Page faults refer to operations that require the database server to access data which isn’t available in active memory. The page_faults counter may increase dramatically during moments of poor performance and may correlate with limited memory environments and larger data sets. Limited and sporadic page faults do not necessarily indicate an issue",
			},
			"heap_usage_bytes": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The heap_usage_bytes field is only available on Unix/Linux systems, and reports the total size in bytes of heap space used by the database process",
			},
		},
		"global_lock": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The globalLock data structure contains information regarding the database’s current lock state, historical lock status, current operation queue, and the number of active clients",
			},
			"ratio": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of ratio displays the relationship between lockTime and totalTime. Low values indicate that operations have held the globalLock frequently for shorter periods of time. High values indicate that operations have held globalLock infrequently for longer periods of time",
			},
			"total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The value of totalTime represents the time, in microseconds, since the database last started and creation of the globalLock. This is roughly equivalent to total server uptime",
			},
			"lock_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The value of lockTime represents the time, in microseconds, since the database last started, that the globalLock has been held",
			},
		},
		"global_lock_current_queue": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"type"},
				Help:   "The currentQueue data structure value provides more granular information concerning the number of operations queued because of a lock",
			},
			"reader": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of readers is the number of operations that are currently queued and waiting for the read lock. A consistently small read-queue, particularly of shorter operations should cause no concern",
			},
			"writer": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of writers is the number of operations that are currently queued and waiting for the write lock. A consistently small write-queue, particularly of shorter operations is no cause for concern",
			},
		},
		"global_lock_client": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"type"},
				Help:   "The activeClients data structure provides more granular information about the number of connected clients and the operation types (e.g. read or write) performed by these clients",
			},
			"reader": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of readers contains a count of the active client connections performing read operations",
			},
			"writer": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of writers contains a count of active client connections performing write operations",
			},
		},
		"index_counters_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type"},
				Help:   "Total indexes by type",
			},
			"accesses": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "accesses reports the number of times that operations have accessed indexes. This value is the combination of the hits and misses. Higher values indicate that your database has indexes and that queries are taking advantage of these indexes. If this number does not grow over time, this might indicate that your indexes do not effectively support your use",
			},
			"hits": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The hits value reflects the number of times that an index has been accessed and mongod is able to return the index from memory. A higher value indicates effective index use. hits values that represent a greater proportion of the accesses value, tend to indicate more effective index configuration",
			},
			"misses": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The misses value represents the number of times that an operation attempted to access an index that was not in memory. These \"misses,\" do not indicate a failed query or operation, but rather an inefficient use of the index. Lower values in this field indicate better index use and likely overall performance as well",
			},
			"resets": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The resets value reflects the number of times that the index counters have been reset since the database last restarted. Typically this value is 0, but use this value to provide context for the data specified by other indexCounters values",
			},
		},
		"index_counters": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The indexCounters data structure reports information regarding the state and use of indexes in MongoDB",
			},
			"miss_ratio": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The missRatio value is the ratio of hits to misses. This value is typically 0 or approaching 0",
			},
		},
		"locks_time_locked_global_microseconds_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type", "database"},
				Help:   "amount of time in microseconds that any database has held the global lock",
			},
			"read": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The R field reports the amount of time in microseconds that any database has held the global read lock",
			},
			"write": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The W field reports the amount of time in microseconds that any database has held the global write lock",
			},
		},
		"locks_time_locked_local_microseconds_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type", "database"},
				Help:   "amount of time in microseconds that any database has held the local lock",
			},
			"read": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The r field reports the amount of time in microseconds that any database has held the local read lock",
			},
			"write": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The w field reports the amount of time in microseconds that any database has held the local write lock",
			},
		},
		"locks_time_acquiring_global_microseconds_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type", "database"},
				Help:   "amount of time in microseconds that any database has spent waiting for the global lock",
			},
			"write": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The W field reports the amount of time in microseconds that any database has spent waiting for the global write lock",
			},
			"read": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The R field reports the amount of time in microseconds that any database has spent waiting for the global read lock",
			},
		},
		"cursors": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"state"},
				Help:   "The cursors data structure contains data regarding cursor state and use",
			},
			"open": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "totalOpen provides the number of cursors that MongoDB is maintaining for clients. Because MongoDB exhausts unused cursors, typically this value small or zero. However, if there is a queue, stale tailable cursor, or a large number of operations, this value may rise.",
			},
			"no_timeout": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "totalNoTimeout provides the number of open cursors with the option DBQuery.Option.noTimeout set to prevent timeout after a period of inactivity.",
			},
			"pinned": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "serverStatus.cursors.pinned provides the number of \"pinned\" open cursors.",
			},
		},
		"cursors_metrics": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The cursors data structure contains data regarding cursor state and use",
			},
			"timed_out_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "timedOut provides a counter of the total number of cursors that have timed out since the server process started. If this number is large or growing at a regular rate, this may indicate an application error.",
			},
		},
		"network_bytes_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"state"},
				Help:   "The network data structure contains data regarding MongoDB’s network use",
			},
			"in_bytes": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The value of the bytesIn field reflects the amount of network traffic, in bytes, received by this database. Use this value to ensure that network traffic sent to the mongod process is consistent with expectations and overall inter-application traffic",
			},
			"out_bytes": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The value of the bytesOut field reflects the amount of network traffic, in bytes, sent from this database. Use this value to ensure that network traffic sent by the mongod process is consistent with expectations and overall inter-application traffic",
			},
		},
		"network_metrics": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The network data structure contains data regarding MongoDB’s network use",
			},
			"num_requests_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "The numRequests field is a counter of the total number of distinct requests that the server has received. Use this value to provide context for the bytesIn and bytesOut values to ensure that MongoDB’s network utilization is consistent with expectations and application use",
			},
		},
		"op_counters_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type"},
				Help:   "The opcounters data structure provides an overview of database operations by type and makes it possible to analyze the load on the database in more granular manner. These numbers will grow over time and in response to database use. Analyze these values over time to track database utilization",
			},
			"insert": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "insert provides a counter of the total number of insert operations received since the mongod instance last started.",
			},
			"query": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "query provides a counter of the total number of queries received since the mongod instance last started",
			},
			"update": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "update provides a counter of the total number of update operations recieved since the mongod instance last started",
			},
			"delete": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "delete provides a counter of the total number of delete operations since the mongod instance last started",
			},
			"getmore": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "getmore provides a counter of the total number of \"getmore\" operations since the mongod instance last started. This counter can be high even if the query count is low. Secondary nodes send getMore operations as part of the replication process",
			},
			"command": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "command provides a counter of the total number of commands issued to the database since the mongod instance last started.",
			},
		},
		"op_counters_repl_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type"},
				Help:   "The opcountersRepl data structure, similar to the opcounters data structure, provides an overview of database replication operations by type and makes it possible to analyze the load on the replica in more granular manner. These values only appear when the current host has replication enabled",
			},
			"insert": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "insert provides a counter of the total number of replicated insert operations since the mongod instance last started",
			},
			"query": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "query provides a counter of the total number of replicated queries since the mongod instance last started",
			},
			"update": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "update provides a counter of the total number of replicated update operations since the mongod instance last started",
			},
			"delete": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "delete provides a counter of the total number of replicated delete operations since the mongod instance last started",
			},
			"getmore": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "getmore provides a counter of the total number of \"getmore\" operations since the mongod instance last started. This counter can be high even if the query count is low. Secondary nodes send getMore operations as part of the replication process",
			},
			"command": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "command provides a counter of the total number of replicated commands issued to the database since the mongod instance last started",
			},
		},
		"memory": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"type"},
				Help:   "The mem data structure holds information regarding the target system architecture of mongod and current memory use",
			},
			"resident": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of resident is roughly equivalent to the amount of RAM, in megabytes (MB), currently used by the database process. In normal use this value tends to grow. In dedicated database servers this number tends to approach the total amount of system memory",
			},
			"virtual": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "virtual displays the quantity, in megabytes (MB), of virtual memory used by the mongod process. With journaling enabled, the value of virtual is at least twice the value of mapped. If virtual value is significantly larger than mapped (e.g. 3 or more times), this may indicate a memory leak",
			},
			"mapped": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "The value of mapped provides the amount of mapped memory, in megabytes (MB), by the database. Because MongoDB uses memory-mapped files, this value is likely to be to be roughly equivalent to the total size of your database or databases",
			},
			"mapped_with_journal": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "mappedWithJournal provides the amount of mapped memory, in megabytes (MB), including the memory used for journaling. This value will always be twice the value of mapped. This field is only included if journaling is enabled",
			},
		},
		"metrics_cursor": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "The cursor is a document that contains data regarding cursor state and use",
			},
			"timed_out_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "timedOut provides the total number of cursors that have timed out since the server process started. If this number is large or growing at a regular rate, this may indicate an application error",
			},
		},
		"metrics_cursor_open": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "gauge_vec",
				Labels: []string{"state"},
				Help:   "The open is an embedded document that contains data regarding open cursors",
			},
			"no_timeout": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "noTimeout provides the number of open cursors with the option DBQuery.Option.noTimeout set to prevent timeout after a period of inactivity",
			},
			"pinned": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "serverStatus.metrics.cursor.open.pinned provides the number of \"pinned\" open cursors",
			},
			"total": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "total provides the number of cursors that MongoDB is maintaining for clients. Because MongoDB exhausts unused cursors, typically this value small or zero. However, if there is a queue, stale tailable cursors, or a large number of operations this value may rise",
			},
		},
		"metrics_document_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"state"},
				Help:   "The document holds a document of that reflect document access and modification patterns and data use. Compare these values to the data in the opcounters document, which track total number of operations",
			},
			"deleted": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "deleted reports the total number of documents deleted",
			},
			"inserted": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "inserted reports the total number of documents inserted",
			},
			"returned": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "returned reports the total number of documents returned by queries",
			},
			"updated": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "updated reports the total number of documents updated",
			},
		},
		"metrics_get_last_error_wtime": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "wtime is a sub-document that reports getLastError operation counts with a w argument greater than 1",
			},
			"num_total": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "num reports the total number of getLastError operations with a specified write concern (i.e. w) that wait for one or more members of a replica set to acknowledge the write operation (i.e. a w value greater than 1.)",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "total_millis reports the total amount of time in milliseconds that the mongod has spent performing getLastError operations with write concern (i.e. w) that wait for one or more members of a replica set to acknowledge the write operation (i.e. a w value greater than 1.)",
			},
		},
		"metrics_get_last_error": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "getLastError is a document that reports on getLastError use",
			},
			"wtimeouts_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "wtimeouts reports the number of times that write concern operations have timed out as a result of the wtimeout threshold to getLastError.",
			},
		},
		"metrics_operation_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type"},
				Help:   "operation is a sub-document that holds counters for several types of update and query operations that MongoDB handles using special operation types",
			},
			"fastmod": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "fastmod reports the number of update operations that neither cause documents to grow nor require updates to the index. For example, this counter would record an update operation that use the $inc operator to increment the value of a field that is not indexed",
			},
			"idhack": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "idhack reports the number of queries that contain the _id field. For these queries, MongoDB will use default index on the _id field and skip all query plan analysis",
			},
			"scan_and_order": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "scanAndOrder reports the total number of queries that return sorted numbers that cannot perform the sort operation using an index",
			},
		},
		"metrics_query_executor_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"state"},
				Help:   "queryExecutor is a document that reports data from the query execution system",
			},
			"scanned": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "scanned reports the total number of index items scanned during queries and query-plan evaluation. This counter is the same as nscanned in the output of explain().",
			},
			"scanned_objects": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "record is a document that reports data related to record allocation in the on-disk memory files",
			},
		},
		"metrics_record": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "record is a document that reports data related to record allocation in the on-disk memory files",
			},
			"moves_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "moves reports the total number of times documents move within the on-disk representation of the MongoDB data set. Documents move as a result of operations that increase the size of the document beyond their allocated record size",
			},
		},
		"metrics_repl_apply_batches": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "batches reports on the oplog application process on secondaries members of replica sets. See Multithreaded Replication for more information on the oplog application processes",
			},
			"num_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "num reports the total number of batches applied across all databases",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "total_millis reports the total amount of time the mongod has spent applying operations from the oplog",
			},
		},
		"metrics_repl_apply": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "apply holds a sub-document that reports on the application of operations from the replication oplog",
			},
			"ops_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "ops reports the total number of oplog operations applied",
			},
		},
		"metrics_repl_buffer": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "MongoDB buffers oplog operations from the replication sync source buffer before applying oplog entries in a batch. buffer provides a way to track the oplog buffer. See Multithreaded Replication for more information on the oplog application process",
			},
			"count": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "count reports the current number of operations in the oplog buffer",
			},
			"max_size_bytes": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "maxSizeBytes reports the maximum size of the buffer. This value is a constant setting in the mongod, and is not configurable",
			},
			"size_bytes": &FieldDesc{
				Type:   "gauge",
				Labels: []string{},
				Help:   "sizeBytes reports the current size of the contents of the oplog buffer",
			},
		},
		"metrics_repl_network_getmores": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "getmores reports on the getmore operations, which are requests for additional results from the oplog cursor as part of the oplog replication process",
			},
			"num_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "num reports the total number of getmore operations, which are operations that request an additional set of operations from the replication sync source.",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "total_millis reports the total amount of time required to collect data from getmore operations",
			},
		},
		"metrics_repl_network": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "network reports network use by the replication process",
			},
			"bytes_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "bytes reports the total amount of data read from the replication sync source",
			},
			"ops_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "ops reports the total number of operations read from the replication source.",
			},
			"readers_created_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "readersCreated reports the total number of oplog query processes created. MongoDB will create a new oplog query any time an error occurs in the connection, including a timeout, or a network operation. Furthermore, readersCreated will increment every time MongoDB selects a new source fore replication.",
			},
		},
		"metrics_repl_oplog_insert": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "insert is a document that reports insert operations into the oplog",
			},
			"num_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "num reports the total number of items inserted into the oplog.",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "total_millis reports the total amount of time spent for the mongod to insert data into the oplog.",
			},
		},
		"metrics_repl_oplog": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "oplog is a document that reports on the size and use of the oplog by this mongod instance",
			},
			"insert_bytes_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "insertBytes the total size of documents inserted into the oplog.",
			},
		},
		"metrics_repl_preload_docs": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "docs is a sub-document that reports on the documents loaded into memory during the pre-fetch stage",
			},
			"num_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "num reports the total number of documents loaded during the pre-fetch stage of replication",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "total_millis reports the total amount of time spent loading documents as part of the pre-fetch stage of replication",
			},
		},
		"metrics_repl_preload_indexes": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "indexes is a sub-document that reports on the index items loaded into memory during the pre-fetch stage of replication",
			},
			"num_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "num reports the total number of index entries loaded by members before updating documents as part of the pre-fetch stage of replication",
			},
			"total_milliseconds": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "total_millis reports the total amount of time spent loading index entries as part of the pre-fetch stage of replication",
			},
		},
		"metrics_storage_freelist_search_total": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "counter_vec",
				Labels: []string{"type"},
				Help:   "metrics about searching records in the database.",
			},
			"bucket_exhausted": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "bucketExhausted reports the number of times that mongod has checked the free list without finding a suitably large record allocation",
			},
			"requests": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "requests reports the number of times mongod has searched for available record allocations",
			},
			"scanned": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "scanned reports the number of available record allocations mongod has searched",
			},
		},
		"metrics_ttl": GroupFieldsMap{
			"metadata": &FieldDesc{
				Type:   "metrics",
				Labels: []string{},
				Help:   "ttl is a sub-document that reports on the operation of the resource use of the ttl index process",
			},
			"deleted_documents_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "deletedDocuments reports the total number of documents deleted from collections with a ttl index.",
			},
			"passes_total": &FieldDesc{
				Type:   "counter",
				Labels: []string{},
				Help:   "passes reports the number of times the background process removes documents from collections with a ttl index",
			},
		},
	}
	// EnabledGroups is map with the group name as field and a boolean indicating wether that group is enabled or not.
	EnabledGroups = make(map[string]bool)
)

// ParseEnabledGroups parses the groups passed by the command line input.
func ParseEnabledGroups(enabledGroupsFlag string) {
	for _, name := range strings.Split(enabledGroupsFlag, ",") {
		name = strings.TrimSpace(name)
		EnabledGroups[name] = true
	}
}

// GroupFields returns a GroupFieldsMap given a groupName.
func GroupFields(groupName string) GroupFieldsMap {
	fields := GroupsDesc[groupName]
	if fields == nil {
		panic("Couldn't find group:" + groupName)
	}

	return fields
}

// GroupField returns a FieldDesc given a group and a field name.
func GroupField(groupName string, fieldName string) *FieldDesc {
	field := GroupFields(groupName)[fieldName]

	if field == nil {
		panic("Couldn't find field: " + fieldName + " in: " + groupName)
	}

	return field
}
