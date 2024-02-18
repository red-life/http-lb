package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
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
	logger, _ := zap.NewDevelopment()
	addrMng := algorithms.NewBackendAddrsManager(backendAddrs, logger)
	rr := algorithms.NewRoundRobin(addrMng, logger)
	for i := 0; i < 100; i++ {
		expected := backendAddrs[i%len(backendAddrs)]
		chosenBackend, err := rr.ChooseBackend(http_lb.Request{})
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if expected != chosenBackend {
			t.Errorf("Failed on %d: Expected '%s' but got '%s'\n", i, expected, chosenBackend)
		}
	}
}
