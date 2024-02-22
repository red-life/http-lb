package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestIPHash_SelectServer(t *testing.T) {
	servers := []string{
		"server 1",
		"server 2",
		"server 3",
		"server 4",
		"server 5",
		"server 6",
	}
	logger, _ := zap.NewDevelopment()
	serverPool := http_lb.NewServerPool(servers, logger)
	ipHash := algorithms.NewIPHash(http_lb.Hash, serverPool, logger)
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{http_lb.Request{RemoteIP: "1.1.1.1"}, servers[int(http_lb.Hash("1.1.1.1"))%len(servers)]},
		{http_lb.Request{RemoteIP: "2.2.2.2"}, servers[int(http_lb.Hash("2.2.2.2"))%len(servers)]},
		{http_lb.Request{RemoteIP: "3.3.3.3"}, servers[int(http_lb.Hash("3.3.3.3"))%len(servers)]},
		{http_lb.Request{RemoteIP: "4.4.4.4"}, servers[int(http_lb.Hash("4.4.4.4"))%len(servers)]},
	}
	for i, test := range tests {
		selectedServer, err := ipHash.SelectServer(test.input)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if test.expected != selectedServer {
			t.Errorf("Failed on %d with IP %s: Expected '%s' but got '%s'\n", i, test.input.RemoteIP,
				test.expected, selectedServer)
		}
	}

}
