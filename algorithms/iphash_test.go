package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
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
	ipHash := algorithms.NewIPHash(backendAddrs, http_lb.Hash)
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{input: http_lb.Request{RemoteIP: "1.1.1.1"}, expected: backendAddrs[int(http_lb.Hash("1.1.1.1"))%len(backendAddrs)]},
		{input: http_lb.Request{RemoteIP: "2.2.2.2"}, expected: backendAddrs[int(http_lb.Hash("2.2.2.2"))%len(backendAddrs)]},
		{input: http_lb.Request{RemoteIP: "3.3.3.3"}, expected: backendAddrs[int(http_lb.Hash("3.3.3.3"))%len(backendAddrs)]},
		{input: http_lb.Request{RemoteIP: "4.4.4.4"}, expected: backendAddrs[int(http_lb.Hash("4.4.4.4"))%len(backendAddrs)]},
	}
	for i, test := range tests {
		chosenBackend := ipHash.ChooseBackend(test.input)
		if test.expected != chosenBackend {
			t.Errorf("Failed on %d with IP %s: Expected '%s' but got '%s'", i, test.input.RemoteIP,
				test.expected, chosenBackend)
		}
	}

}
