package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/dcu/mongodb_exporter/collector"
    "github.com/dcu/mongodb_exporter/shared"
    "flag"
    //"github.com/golang/glog"
)
var (
	listenAddress     = flag.String("web.listen-address", ":9001", "Address on which to expose metrics and web interface.")
	metricsPath       = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")

  mongodbUri       = flag.String("mongodb.uri", "mongodb://localhost:27017", "Mongodb URI, format: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]")
	//enabledCollectors = flag.String("collectors.enabled", "instance,asserts,background_flushing,connections,durability,durability_time_ms,extra_info,global_lock,global_lock_current_queue,global_lock_client,index_counters,locks_time_locked_micros,locks_time_acquiring_micros,network,op_counters,op_counters_repl,memory,metrics_cursor_timed_out,metrics_cursor_open,metrics_document,metrics_get_last_error_wtime,metrics_get_last_error,metrics_operation,metrics_query_executor,metrics_record,metrics_repl_apply_batches,metrics_repl_apply,metrics_repl_buffer,metrics_repl_network_getmores,metrics_repl_network,metrics_repl_oplog_insert,metrics_repl_oplog,metrics_repl_preload_docs,metrics_repl_preload_indexes,metrics_storage,metrics_ttl", "Comma-separated list of collectors to use.")
	//printCollectors   = flag.Bool("collectors.print", false, "If true, print available collectors and exit.")
	authUser          = flag.String("auth.user", "", "Username for basic auth.")
	authPass          = flag.String("auth.pass", "", "Password for basic auth.")
)

type basicAuthHandler struct {
	handler  http.HandlerFunc
	user     string
	password string
}

func (h *basicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok || password != h.password || user != h.user {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"metrics\"")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	h.handler(w, r)
	return
}

func hasUserAndPassword() bool {
    return *authUser != "" && *authPass != ""
}

func prometheusHandler() http.Handler {
    handler := prometheus.Handler()
    if hasUserAndPassword() {
        handler = &basicAuthHandler{
            handler:  prometheus.Handler().ServeHTTP,
            user:     *authUser,
            password: *authPass,
        }
    }

    return handler
}

func startWebServer() {
    handler := prometheusHandler()

    http.Handle(*metricsPath, handler)
    err := http.ListenAndServe(":9001", nil)

    if err != nil {
        panic(err)
    }

}

func main() {
    flag.Parse()
    shared.LoadGroupsDesc()

    mongodbCollector := collector.NewMongodbCollector(collector.MongodbCollectorOpts{
        URI: *mongodbUri,
    })
    prometheus.MustRegister(mongodbCollector)

    startWebServer()
}

