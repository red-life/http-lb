package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"go.uber.org/zap"
	"testing"
)

func TestIPHash_ChooseBackend(t *testing.T) {
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
	ipHash := algorithms.NewIPHash(http_lb.Hash, addrMng, logger)
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{http_lb.Request{RemoteIP: "1.1.1.1"}, backendAddrs[int(http_lb.Hash("1.1.1.1"))%len(backendAddrs)]},
		{http_lb.Request{RemoteIP: "2.2.2.2"}, backendAddrs[int(http_lb.Hash("2.2.2.2"))%len(backendAddrs)]},
		{http_lb.Request{RemoteIP: "3.3.3.3"}, backendAddrs[int(http_lb.Hash("3.3.3.3"))%len(backendAddrs)]},
		{http_lb.Request{RemoteIP: "4.4.4.4"}, backendAddrs[int(http_lb.Hash("4.4.4.4"))%len(backendAddrs)]},
	}
	for i, test := range tests {
		chosenBackend, err := ipHash.ChooseBackend(test.input)
		if err != nil {
			t.Fatalf("Expected err to be nil but got %s\n", err)
		}
		if test.expected != chosenBackend {
			t.Errorf("Failed on %d with IP %s: Expected '%s' but got '%s'\n", i, test.input.RemoteIP,
				test.expected, chosenBackend)
		}
	}

}
