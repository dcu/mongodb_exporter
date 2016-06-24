package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	slog "log"
	"net/http"
	"os"
	"strings"

	"github.com/dcu/mongodb_exporter/collector"
	"github.com/dcu/mongodb_exporter/shared"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

func mongodbDefaultUri() string {
	if u := os.Getenv("MONGODB_URL"); u != "" {
		return u
	}
	return "mongodb://localhost:27017"
}

var (
	listenAddressFlag = flag.String("web.listen-address", ":9001", "Address on which to expose metrics and web interface.")
	metricsPathFlag   = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
	webTlsCert        = flag.String("web.tls-cert", "", "Path to PEM file that conains the certificate (and opionally also the private key in PEM format).\n"+
		"    \tThis should include the whole certificate chain.\n"+
		"    \tIf provided: The web socket will be a HTTPS socket.\n"+
		"    \tIf not provided: Only HTTP.")
	webTlsPrivateKey = flag.String("web.tls-private-key", "", "Path to PEM file that conains the private key (if not contained in web.tls-cert file).")
	webTlsClientCa   = flag.String("web.tls-client-ca", "", "Path to PEM file that conains the CAs that are trused for client connections.\n"+
		"    \tIf provided: Connecting clients should present a certificate signed by one of this CAs.\n"+
		"    \tIf not provided: Every client will be accepted.")

	mongodbURIFlag = flag.String("mongodb.uri", mongodbDefaultUri(), "Mongodb URI, format: [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]")
	mongodbTlsCert = flag.String("mongodb.tls-cert", "", "Path to PEM file that conains the certificate (and opionally also the private key in PEM format).\n"+
		"    \tThis should include the whole certificate chain.\n"+
		"    \tIf provided: The connection will be opened via TLS to the MongoDB server.")
	mongodbTlsPrivateKey = flag.String("mongodb.tls-private-key", "", "Path to PEM file that conains the private key (if not contained in mongodb.tls-cert file).")
	mongodbTlsCa         = flag.String("mongodb.tls-ca", "", "Path to PEM file that conains the CAs that are trused for server connections.\n"+
		"    \tIf provided: MongoDB servers connecting to should present a certificate signed by one of this CAs.\n"+
		"    \tIf not provided: System default CAs are used.")
	mongodbTlsDisableHostnameValidation = flag.Bool("mongodb.tls-disable-hostname-validation", false, "Do hostname validation for server connection.")
	enabledGroupsFlag                   = flag.String("groups.enabled", "asserts,durability,background_flushing,connections,extra_info,global_lock,index_counters,network,op_counters,op_counters_repl,memory,locks,metrics", "Comma-separated list of groups to use, for more info see: docs.mongodb.org/manual/reference/command/serverStatus/")
	authUserFlag                        = flag.String("auth.user", "", "Username for basic auth.")
	authPassFlag                        = flag.String("auth.pass", "", "Password for basic auth.")
	mongodbCollectOplog                        = flag.Bool("mongodb.collect.oplog", true, "collect Mongodb Oplog status")
	mongodbCollectReplSet                      = flag.Bool("mongodb.collect.replset", true, "collect Mongodb replica set status")
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
	return *authUserFlag != "" && *authPassFlag != ""
}

func prometheusHandler() http.Handler {
	handler := prometheus.Handler()
	if hasUserAndPassword() {
		handler = &basicAuthHandler{
			handler:  prometheus.Handler().ServeHTTP,
			user:     *authUserFlag,
			password: *authPassFlag,
		}
	}

	return handler
}

func startWebServer() {
	handler := prometheusHandler()

	registerCollector()

	http.Handle(*metricsPathFlag, handler)

	server := &http.Server{
		Addr:     *listenAddressFlag,
		ErrorLog: createHttpServerLogWrapper(),
	}

	var err error
	if len(*webTlsCert) > 0 {
		clientValidation := "no"
		if len(*webTlsClientCa) > 0 && len(*webTlsCert) > 0 {
			certificates, err := shared.LoadCertificatesFrom(*webTlsClientCa)
			if err != nil {
				glog.Fatalf("Couldn't load client CAs from %s. Got: %s", *webTlsClientCa, err)
			}
			server.TLSConfig = &tls.Config{
				ClientCAs:  certificates,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
			clientValidation = "yes"
		}
		targetTlsPrivateKey := *webTlsPrivateKey
		if len(targetTlsPrivateKey) <= 0 {
			targetTlsPrivateKey = *webTlsCert
		}
		fmt.Printf("Listening on %s (scheme=HTTPS, secured=TLS, clientValidation=%s)\n", server.Addr, clientValidation)
		err = server.ListenAndServeTLS(*webTlsCert, targetTlsPrivateKey)
	} else {
		fmt.Printf("Listening on %s (scheme=HTTP, secured=no, clientValidation=no)\n", server.Addr)
		err = server.ListenAndServe()
	}

	if err != nil {
		panic(err)
	}
}

func registerCollector() {
	mongodbCollector := collector.NewMongodbCollector(collector.MongodbCollectorOpts{
		URI:                   *mongodbURIFlag,
		TLSCertificateFile:    *mongodbTlsCert,
		TLSPrivateKeyFile:     *mongodbTlsPrivateKey,
		TLSCaFile:             *mongodbTlsCa,
		TLSHostnameValidation: !(*mongodbTlsDisableHostnameValidation),
		CollectOplog: *mongodbCollectOplog,
		CollectReplSet: *mongodbCollectReplSet,
	})
	prometheus.MustRegister(mongodbCollector)
}

type bufferedLogWriter struct {
	buf []byte
}

func (w *bufferedLogWriter) Write(p []byte) (n int, err error) {
	glog.Info(strings.TrimSpace(strings.Replace(string(p), "\n", " ", -1)))
	return len(p), nil
}

func createHttpServerLogWrapper() *slog.Logger {
	return slog.New(&bufferedLogWriter{}, "", 0)
}

func main() {
	flag.Parse()
	shared.ParseEnabledGroups(*enabledGroupsFlag)

	startWebServer()
}
