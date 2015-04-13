package collector

import (
	"github.com/dcu/mongodb_exporter/shared"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	shared.LoadGroupsDesc()
	os.Exit(m.Run())
}

func LoadFixture(name string) []byte {
	data, err := ioutil.ReadFile("fixtures/" + name)
	if err != nil {
		panic(err)
	}

	return data
}
