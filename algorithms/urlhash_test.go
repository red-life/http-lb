package algorithms_test

import (
	"github.com/red-life/http-lb"
	"github.com/red-life/http-lb/algorithms"
	"testing"
)

func TestURLHash_ChooseBackend(t *testing.T) {
	backendAddrs := []string{
		"addr 1",
		"addr 2",
		"addr 3",
		"addr 4",
		"addr 5",
		"addr 6",
	}
	urlHash := algorithms.NewURLHash(http_lb.Hash, algorithms.NewBackendAddrsManager(backendAddrs))
	tests := []struct {
		input    http_lb.Request
		expected string
	}{
		{input: http_lb.Request{URLPath: "/"}, expected: backendAddrs[int(http_lb.Hash("/"))%len(backendAddrs)]},
		{input: http_lb.Request{URLPath: "/home"}, expected: backendAddrs[int(http_lb.Hash("/home"))%len(backendAddrs)]},
		{input: http_lb.Request{URLPath: "/auth/login"}, expected: backendAddrs[int(http_lb.Hash("/auth/login"))%len(backendAddrs)]},
		{input: http_lb.Request{URLPath: "/api/v1"}, expected: backendAddrs[int(http_lb.Hash("/api/v1"))%len(backendAddrs)]},
	}
	for i, test := range tests {
		chosenBackend := urlHash.ChooseBackend(test.input)
		if test.expected != chosenBackend {
			t.Errorf("Failed on %d with Path %s: Expected '%s' but got '%s'", i, test.input.URLPath,
				test.expected, chosenBackend)
		}
	}

}
