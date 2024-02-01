package algorithms_test

import (
	"github.com/red-life/http-lb/algorithms"
	"testing"
)

func TestRoundRobin_ChooseBackend(t *testing.T) {
	backendAddrs := []string{
		"addr 1",
		"addr 2",
		"addr 3",
		"addr 4",
		"addr 5",
		"addr 6",
	}
	rr := algorithms.NewRoundRobin(backendAddrs)
	for i := 0; i < 100; i++ {
		expected := backendAddrs[i%len(backendAddrs)]
		chosenBackend := rr.ChooseBackend()
		if expected != chosenBackend {
			t.Errorf("Failed on %d: Expected '%s' but got '%s'", i, expected, chosenBackend)
		}
	}
}