package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestStickyRoundRobin_SelectServer(t *testing.T) {
	servers := []string{
		"server 1",
		"server 2",
		"server 3",
		"server 4",
		"server 5",
		"server 6",
	}
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{http_lb.Request{RemoteIP: "1.1.1.1"}, servers[0]},
		{http_lb.Request{RemoteIP: "2.2.2.2"}, servers[1]},
		{http_lb.Request{RemoteIP: "1.1.1.1"}, servers[0]},
		{http_lb.Request{RemoteIP: "3.3.3.3"}, servers[2]},
		{http_lb.Request{RemoteIP: "4.4.4.4"}, servers[3]},
		{http_lb.Request{RemoteIP: "5.5.5.5"}, servers[4]},
		{http_lb.Request{RemoteIP: "6.6.6.6"}, servers[5]},
		{http_lb.Request{RemoteIP: "3.3.3.3"}, servers[2]},
		{http_lb.Request{RemoteIP: "4.4.4.4"}, servers[3]},
		{http_lb.Request{RemoteIP: "6.6.6.6"}, servers[5]},
	}
	logger, _ := zap.NewDevelopment()
	serverPool := http_lb.NewServerPool(servers, logger)
	sticky_rr := algorithms.NewStickyRoundRobin(serverPool, logger)
	for i, test := range tests {
		selectedServer, err := sticky_rr.SelectServer(test.input)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if selectedServer != test.expected {
			t.Errorf("Failed on %d with IP %s: Expected '%s' but got '%s'\n", i, test.input.RemoteIP,
				test.expected, selectedServer)
		}
	}
}
