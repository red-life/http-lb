package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestStickyRoundRobin_ChooseBackend(t *testing.T) {
	backendAddrs := []string{
		"addr 1",
		"addr 2",
		"addr 3",
		"addr 4",
		"addr 5",
		"addr 6",
	}
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{http_lb.Request{RemoteIP: "1.1.1.1"}, backendAddrs[0]},
		{http_lb.Request{RemoteIP: "2.2.2.2"}, backendAddrs[1]},
		{http_lb.Request{RemoteIP: "1.1.1.1"}, backendAddrs[0]},
		{http_lb.Request{RemoteIP: "3.3.3.3"}, backendAddrs[2]},
		{http_lb.Request{RemoteIP: "4.4.4.4"}, backendAddrs[3]},
		{http_lb.Request{RemoteIP: "5.5.5.5"}, backendAddrs[4]},
		{http_lb.Request{RemoteIP: "6.6.6.6"}, backendAddrs[5]},
		{http_lb.Request{RemoteIP: "3.3.3.3"}, backendAddrs[2]},
		{http_lb.Request{RemoteIP: "4.4.4.4"}, backendAddrs[3]},
		{http_lb.Request{RemoteIP: "6.6.6.6"}, backendAddrs[5]},
	}
	logger, _ := zap.NewDevelopment()
	addrMng := algorithms.NewBackendAddrsManager(backendAddrs, logger)
	sticky_rr := algorithms.NewStickyRoundRobin(addrMng, logger)
	for i, test := range tests {
		chosenBackend, err := sticky_rr.ChooseBackend(test.input)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if chosenBackend != test.expected {
			t.Errorf("Failed on %d with IP %s: Expected '%s' but got '%s'\n", i, test.input.RemoteIP,
				test.expected, chosenBackend)
		}
	}
}
