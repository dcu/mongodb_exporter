test:
	go test github.com/dcu/mongodb_exporter/collector -cover -coverprofile=collector_coverage.out -short
	go tool cover -func=collector_coverage.out
	go test github.com/dcu/mongodb_exporter/shared -cover -coverprofile=shared_coverage.out -short
	go tool cover -func=shared_coverage.out
	@rm *.out

deps:
	go get -u github.com/prometheus/client_golang/prometheus
	go get -u gopkg.in/mgo.v2

build: deps
	go build mongodb_exporter.go

