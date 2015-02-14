test:
	go test github.com/dcu/mongodb_exporter/collector -cover -coverprofile=collector_coverage.out -short
	go tool cover -func=collector_coverage.out
	go test github.com/dcu/mongodb_exporter/shared -cover -coverprofile=shared_coverage.out -short
	go tool cover -func=shared_coverage.out
	@rm *.out

bindata:
	go-bindata -pkg=shared -o=shared/assets.go groups.yml

deps:
	go get -u github.com/prometheus/client_golang/prometheus
	go get -u gopkg.in/yaml.v2
	go get -u github.com/jteeuwen/go-bindata
	go get -u gopkg.in/mgo.v2

build: deps bindata
	go build mongodb_exporter.go


