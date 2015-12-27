package collector

import (
	"testing"
)

func Test_MemoryCollectData(t *testing.T) {
	stats := &MemStats{}

	stats.Export()
}
