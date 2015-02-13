package collector

import(
	"testing"
)

func Test_CollectServerStatus(t *testing.T) {
	collector := NewMongodbCollector(MongodbCollectorOpts{URI: "localhost"})
	serverStatus := collector.collectServerStatus(nil)

	if serverStatus.Asserts == nil {
		t.Error("Error loading document.")
	}
}


