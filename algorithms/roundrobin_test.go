package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestRoundRobin_SelectServer(t *testing.T) {
	servers := []string{
		"server 1",
		"server 2",
		"server 3",
		"server 4",
		"server 5",
		"server 6",
	}
	logger, _ := zap.NewDevelopment()
	serverPool := algorithms.NewServerPool(servers, logger)
	rr := algorithms.NewRoundRobin(serverPool, logger)
	for i := 0; i < 100; i++ {
		expected := servers[i%len(servers)]
		selectedServer, err := rr.SelectServer(http_lb.Request{})
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if expected != selectedServer {
			t.Errorf("Failed on %d: Expected '%s' but got '%s'\n", i, expected, selectedServer)
		}
	}
}
