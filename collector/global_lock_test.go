package collector

import (
	"testing"
)

func Test_GlobalLockCollectData(t *testing.T) {
	stats := &GlobalLockStats{
		CurrentQueue:  &QueueStats{},
		ActiveClients: &ClientStats{},
	}

	stats.Export()
}
